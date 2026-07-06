package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/nhassl3/hairdress_arz/internal/repository/redis"
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
	verifyRedis    domain.VerifyRedis
	sender         domain.VerifySender
}

type TokenPair struct {
	AccessToken,
	RefreshToken string
	RefreshTokenPayload *authServiceHub.Payload
}

func NewAuthService(userRepo domain.UserRepository, accessManager, refreshManager auth.TokenManager,
	blacklist auth.TokenBlacklist, userRedis domain.UserRedis, verifyRedis domain.VerifyRedis, sender domain.VerifySender,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		refreshManager: refreshManager,
		accessManager:  accessManager,
		blacklist:      blacklist,
		userRedis:      userRedis,
		verifyRedis:    verifyRedis,
		sender:         sender,
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

func (s *AuthService) RequestVerifyEmail(ctx context.Context, email string, operationId string) (string, error) {
	if exists, err := s.userRepo.ExistsByUsername(ctx, email); !exists {
		if err != nil {
			return "", fmt.Errorf("auth_service.RequestVerifyEmail: failed to check existence of user by username: %w", err)
		}
		return "", domain.ErrUserNotFound
	}

	// If operationId is provided, this is a re-send — check existing record
	if operationId != "" {
		record, err := s.verifyRedis.Code(ctx, operationId)
		if err != nil {
			return "", fmt.Errorf("auth_service.RequestVerifyEmail: failed to get existing code: %w", err)
		}
		if record.AttemptsLeft <= 0 {
			return "", domain.ErrSmsRateLimited
		}
		// Delete old record so SetCode with NX can succeed
		_ = s.verifyRedis.DelCode(ctx, operationId)
	}

	var code string
	for range 5 {
		code = s.sender.Helper().GenerateVerifyCode()
		if code != "" {
			break
		}
	}
	if code == "" {
		return "", fmt.Errorf("auth_service.RequestVerifyEmail: failed to generate verify code")
	}

	if operationId == "" {
		operationId = uuid.NewString()
	}

	if err := s.verifyRedis.SetCode(ctx, operationId, domain.NewVerifyState(
		&domain.MethodToVerify{Email: &email}, s.sender.Helper().Code2Hash(code), 3, 0)); err != nil {
		return "", fmt.Errorf("auth_service.RequestVerifyEmail: failed to set verify code in Redis: %w", err)
	}

	if err := s.sender.SendEmail(code, email); err != nil {
		return "", fmt.Errorf("auth_service.RequestVerifyEmail: failed to send verify code: %w", err)
	}

	return operationId, nil
}

func (s *AuthService) ApproveCode(ctx context.Context, operationId string, method domain.MethodToVerify, code string) (string, error) {
	record, err := s.verifyRedis.Code(ctx, operationId)
	if err != nil {
		return "", fmt.Errorf("auth_service.ApproveCode: failed to get code: %w", err)
	}
	if time.Now().After(record.ExpiresAt) {
		return "", fmt.Errorf("auth_service.ApproveCode: expired code")
	}
	if record.AttemptsLeft <= 0 {
		return "", domain.ErrSmsRateLimited
	}

	if ok := s.sender.Helper().CompareCode(code, record.HashCode); !ok {
		return "", domain.ErrInvalidCode
	}

	_ = s.verifyRedis.DelCode(ctx, operationId)

	token := uuid.NewString()

	entryCode := redis.VerifySmsEntryKey
	if method.Email != nil && *method.Email != "" {
		entryCode = redis.VerifyEmailEntryKey
	}

	if err := s.verifyRedis.SetVerified(ctx, entryCode, token, method); err != nil {
		return "", fmt.Errorf("auth_service.ApproveCode: failed to save verified method: %w", err)
	}

	return token, nil
}

func (s *AuthService) verifyPhoneNumber(ctx context.Context, code, phoneNumber string) error {
	if exists, err := s.userRepo.ExistsByPhoneNumber(ctx, phoneNumber); !exists {
		if err != nil {
			return fmt.Errorf("auth_service.verifyPhoneNumber: failed to check existence of user by phone: %w", err)
		}
		return domain.ErrUserNotFound
	}

	record, err := s.verifyRedis.Code(ctx, redis.VerifySmsEntryKey+phoneNumber)
	if err != nil {
		return fmt.Errorf("auth_service.verifyPhoneNumber: failed to get sms code: %w", err)
	}
	if time.Now().After(record.ExpiresAt) {
		return fmt.Errorf("auth_service.verifyPhoneNumber: expired sms code")
	}
	if record.AttemptsLeft <= 0 {
		return domain.ErrSmsRateLimited
	}

	if ok := s.sender.Helper().CompareCode(code, record.HashCode); !ok {
		return domain.ErrInvalidCode
	}

	_ = s.verifyRedis.DelCode(ctx, redis.VerifySmsEntryKey+phoneNumber)

	return nil
}

func (s *AuthService) Logout(ctx context.Context, payload *authServiceHub.Payload) error {
	if payload != nil && payload.JTI != "" {
		if err := s.blacklist.Blacklist(ctx, payload.JTI, payload.ExpiredAt); err != nil {
			return fmt.Errorf("auth_service.Logout: blacklist access token: %w", err)
		}
		if err := s.userRedis.DelProfile(ctx, payload.Username); err != nil {
			return fmt.Errorf("auth_service.Logout: failed to delete profile: %w", err)
		}
		if err := s.userRedis.DelSession(ctx, payload.Username); err != nil {
			return fmt.Errorf("auth_service.Logout: failed to delete user session: %w", err)
		}
		return s.userRepo.DeleteSession(ctx, payload.Username)
	}
	return domain.ErrInvalidToken
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	payload, err := s.refreshManager.VerifyToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	session, err := s.userRedis.Session(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, domain.ErrRedisNotFound) {
			session, err = s.userRepo.GetSession(ctx, refreshToken)
			if err != nil {
				return nil, fmt.Errorf("auth_service.RefreshToken get session: %w", err)
			}
			if err := s.userRedis.SetSession(ctx, session); err != nil {
				return nil, fmt.Errorf("auth_service.RefreshToken set session (redis): %w", err)
			}
		} else {
			return nil, fmt.Errorf("auth_service.RefreshToken get session (redis): %w", err)
		}
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, domain.ErrExpiredToken
	} else if session.IsBlocked {
		return nil, domain.ErrSessionIsBlocked
	}

	accessToken, err := s.accessManager.CreateToken(payload.Username, payload.UID, payload.Role)
	if err != nil {
		return nil, fmt.Errorf("auth_service: create access token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.userRedis.Profile(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			user, err = s.userRepo.GetByUsername(ctx, username)
			if err != nil {
				return nil, fmt.Errorf("auth_service.GetMe: failed to get user by username: %w", err)
			}
		}
		return nil, fmt.Errorf("auth_service.GetMe: failed to get user by username: %w", err)
	}
	return user, nil
}

// helpers

func generateNewUsername() string {
	return fmt.Sprintf("@user_%d", uuid.New().ID())
}

func (s *AuthService) createSendCodeWebhook(ctx context.Context, phoneNumber string) error {
	code := s.sender.Helper().GenerateVerifyCode()
	if err := s.sender.SendPhone(phoneNumber, code); err != nil {
		return fmt.Errorf("auth_service.createSendCodeWebhook: failed to send sms: %w", err)
	}

	if err := s.verifyRedis.SetCode(
		ctx,
		redis.VerifySmsEntryKey+phoneNumber,
		domain.NewVerifyState(&domain.MethodToVerify{PhoneNumber: &phoneNumber}, s.sender.Helper().Code2Hash(code), 3, 0),
	); err != nil {
		return fmt.Errorf("auth_service.createSendCodeWebhook: failed to save sms verification code (redis): %w", err)
	}

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context) error {
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
