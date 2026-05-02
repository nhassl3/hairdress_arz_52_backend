package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/nhassl3/servicehub-backend/pkg/auth"
	authServiceHub "github.com/nhassl3/servicehub-backend/pkg/auth"
	"google.golang.org/grpc/metadata"
)

type AuthService struct {
	userRepo       domain.UserRepository
	refreshManager auth.TokenManager
	accessManager  auth.TokenManager
	blacklist      auth.TokenBlacklist
	userRedis      domain.UserRedis
	smsRedis       domain.SmsRedis
	smsSender      domain.SmsSender
}

type TokenPair struct {
	AccessToken,
	RefreshToken string
	RefreshTokenPayload *authServiceHub.Payload
}

func NewAuthService(userRepo domain.UserRepository, accessManager, refreshManager auth.TokenManager,
	blacklist auth.TokenBlacklist, userRedis domain.UserRedis, smsRedis domain.SmsRedis, smsSender domain.SmsSender) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		refreshManager: refreshManager,
		accessManager:  accessManager,
		blacklist:      blacklist,
		userRedis:      userRedis,
		smsRedis:       smsRedis,
		smsSender:      smsSender,
	}
}

func (s *AuthService) Login(ctx context.Context, phoneNumber string) (*domain.User, error) {
	clientIP, _ := getMetadataFromContext(ctx) // cooldown 5 minutes for one IP address

	if block, ttl, err := s.userRedis.AuthBlock(ctx, clientIP); block || err != nil {
		return nil, fmt.Errorf("%w: %f", domain.ErrAuthBlock, ttl)
	}

	user, err := s.userRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("auth_service.Login: failed to get user by phone number: %w", err)
	}
	if err := s.userRedis.SetProfile(ctx, user.Username, user); err != nil {
		return nil, fmt.Errorf("auth_service.Register: failed to set profile in Redis: %w", err)
	}

	if err := s.createSendCodeWebhook(ctx, phoneNumber); err != nil {
		return nil, fmt.Errorf("auth_service.Login: failed to create send code webhook: %w", err)
	}

	if err := s.userRedis.SetAuthBlock(ctx, clientIP); err != nil {
		return nil, fmt.Errorf("auth_service.Login: failed to set auth block in Redis: %w", err)
	}

	return user, nil
}

func (s *AuthService) Register(ctx context.Context, params domain.CreateUserParams) (*domain.User, error) {
	clientIP, _ := getMetadataFromContext(ctx) // cooldown 5 minutes for one IP address

	if block, ttl, err := s.userRedis.AuthBlock(ctx, clientIP); block || err != nil {
		return nil, fmt.Errorf("%w: %f", domain.ErrAuthBlock, ttl)
	}

	if params.Username == nil {
		v := generateNewUsername()
		existsByUsername, err := s.userRepo.ExistsByUsername(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("auth_service.Register: failed to check new username: %w", err)
		}
		for existsByUsername {
			v = generateNewUsername()
			existsByUsername, err = s.userRepo.ExistsByUsername(ctx, v)
			if err != nil {
				return nil, fmt.Errorf("auth_service.Register: failed to check new username: %w", err)
			}
		}
		params.Username = &v
	}
	user, err := s.userRepo.Create(ctx, &params)
	if err != nil {
		return nil, err
	}
	if err := s.userRedis.SetProfile(ctx, user.Username, user); err != nil {
		return nil, fmt.Errorf("auth_service.Register: failed to set profile in Redis: %w", err)
	}

	if err := s.createSendCodeWebhook(ctx, params.PhoneNumber); err != nil {
		return nil, err
	}

	if err := s.userRedis.SetAuthBlock(ctx, clientIP); err != nil {
		return nil, fmt.Errorf("auth_service.Login: failed to set auth block in Redis: %w", err)
	}

	return user, nil
}

func (s *AuthService) VerifyCode(ctx context.Context, phone, code string) (*TokenPair, error) {
	user, err := s.userRepo.GetByPhoneNumber(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("auth_service.VerifyCode: failed to get user by phone: %w", err)
	}

	smsRecorder, err := s.smsRedis.GetCode(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("auth_service.VerifyCode: failed to get sms code: %w", err)
	}
	if time.Now().After(smsRecorder.ExpiresAt) {
		return nil, fmt.Errorf("auth_service.VerifyCode: expired sms code")
	}
	if smsRecorder.AttemptsLeft <= 0 {
		return nil, domain.ErrSmsRateLimited
	}

	if ok := s.smsSender.Helper().CompareCode(code, smsRecorder.Hash); !ok {
		return nil, domain.ErrInvalidCode
	}

	tokenPair, err := s.createTokenPair(user.Username, user.UID, "user")
	if err != nil {
		return nil, fmt.Errorf("auth_service.VerifyCode: failed to create token pair: %w", err)
	}

	if err := s.smsRedis.DeleteCode(ctx, phone); err != nil {
		return nil, fmt.Errorf("auth_service.VerifyCode: failed to delete sms code (redis): %w", err)
	}

	return tokenPair, nil
}

func (s *AuthService) Logout(ctx context.Context, payload *authServiceHub.Payload) error {
	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	return nil, nil
}

func (s *AuthService) GetMe(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("auth_service.GetMe: failed to get user by username: %w", err)
	}
	return user, nil
}

// helpers

func generateNewUsername() string {
	return fmt.Sprintf("@user_%d", uuid.New().ID())
}

func (s *AuthService) createSendCodeWebhook(ctx context.Context, phoneNumber string) error {
	code, err := s.smsSender.Helper().GenerateSecureCode()
	if err != nil {
		return fmt.Errorf("auth_service.createSendCodeWebhook: failed to generate secure code: %w", err)
	}
	if err := s.smsSender.Send(ctx, phoneNumber, code); err != nil {
		return fmt.Errorf("auth_service.createSendCodeWebhook: failed to send sms: %w", err)
	}

	if _, err := s.smsRedis.SaveCode(
		ctx,
		phoneNumber,
		s.smsSender.Helper().Code2Hash(code),
	); err != nil {
		return fmt.Errorf("auth_service.createSendCodeWebhook: failed to save sms verification code (redis): %w", err)
	}

	return nil
}

// createTokenPair
func (s *AuthService) createTokenPair(username, uid, role string) (*TokenPair, error) {
	accessToken, err := s.accessManager.CreateToken(username, uid, role)
	if err != nil {
		return nil, fmt.Errorf("auth_service: create access token: %w", err)
	}

	refreshToken, payload, err := s.refreshManager.CreateRefreshToken(username, uid, role)
	if err != nil {
		return nil, fmt.Errorf("auth_service: create refresh token: %w", err)
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, RefreshTokenPayload: payload}, nil
}

// getMetadataFromContext
func getMetadataFromContext(ctx context.Context) (clientIp string, userAgent string) {
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		xForwardFor := headers.Get("x-forwarded-for")
		if len(xForwardFor) > 0 && xForwardFor[0] != "" {
			ips := strings.Split(xForwardFor[0], ",")
			if len(ips) > 0 {
				clientIp = ips[0]
			}
		}
		usrAgent := headers.Get("user-agent")
		if len(usrAgent) >= 1 && usrAgent[0] != "" {
			userAgent = usrAgent[0]
		}
	}
	return
}
