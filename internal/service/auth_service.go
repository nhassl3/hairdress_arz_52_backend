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
	"go.uber.org/zap"
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

func (s *AuthService) Login(ctx context.Context, params *domain.LoginParams) (operationId string, user *domain.User, err error) {
	clientIP, _ := getMetadataFromContext(ctx) // cooldown 5 minutes for one IP address

	zap.L().Info("ClientIP", zap.String("clientIP", clientIP))

	if block, ttl, e := s.userRedis.AuthBlock(ctx, clientIP); block || e != nil {
		return "", nil, fmt.Errorf("%w: %f", domain.ErrAuthBlock, ttl)
	}

	var identifierType string

	switch {
	case params.Username != nil:
		user, err = s.userRepo.GetByUsername(ctx, *params.Username)
		identifierType = "username"
	case params.Email != nil:
		user, err = s.userRepo.GetByEmail(ctx, *params.Email)
		identifierType = "email"
	case params.PhoneNumber != nil:
		user, err = s.userRepo.GetByPhoneNumber(ctx, *params.PhoneNumber)
		identifierType = "phone number"
	default:
		return "", nil, fmt.Errorf("auth_service.Login: no login identifier provided")
	}
	if err != nil {
		return "", nil, fmt.Errorf("auth_service.Login: failed to get user by %s: %w", identifierType, err)
	}

	if err := s.userRedis.SetProfile(ctx, user.Username, user); err != nil {
		return "", nil, fmt.Errorf("auth_service.Login: failed to set profile in Redis: %w", err)
	}

	switch identifierType {
	case "phone number":
		operationId, err = s.createSendCodeWebhook(ctx, *params.PhoneNumber, clientIP)
	case "email":
		operationId, err = s.createSendCodeEmail(ctx, user.Email, clientIP)
	default: // username — send code to the user's registered contact
		if user.PhoneNumber != "" {
			operationId, err = s.createSendCodeWebhook(ctx, user.PhoneNumber, clientIP)
		} else if user.Email != "" {
			operationId, err = s.createSendCodeEmail(ctx, user.Email, clientIP)
		}
	}
	if err != nil {
		return "", nil, fmt.Errorf("auth_service.Login: failed to send verification code: %w", err)
	}

	if err := s.userRedis.SetAuthBlock(ctx, clientIP); err != nil {
		return "", nil, fmt.Errorf("auth_service.Login: failed to set auth block in Redis: %w", err)
	}

	return operationId, user, nil
}

func (s *AuthService) createSendCodeEmail(ctx context.Context, email, clientIP string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("auth_service.createSendCodeEmail: email is empty")
	}

	if err := s.verifyRedis.IncDailyByMethod(ctx, &domain.MethodToVerify{Email: &email}); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeEmail: daily limit exceeded: %w", err)
	}
	if err := s.verifyRedis.IncDailyByIP(ctx, redis.VerifyEmailEntryKey, clientIP); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeEmail: daily IP limit exceeded: %w", err)
	}

	code := s.sender.Helper().GenerateVerifyCode()

	operationId := redis.VerifyEmailEntryKey + email
	if err := s.verifyRedis.SetCode(
		ctx,
		operationId,
		domain.NewVerifyState(&domain.MethodToVerify{Email: &email}, s.sender.Helper().Code2Hash(code), s.verifyRedis.Attempts(), s.verifyRedis.CodeTTL()),
	); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeEmail: failed to save email verification code (redis): %w", err)
	}

	if err := s.sender.SendEmail(email, code); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeEmail: failed to send email: %w", err)
	}

	return operationId, nil
}

func (s *AuthService) LoginVerify(ctx context.Context, verifyToken string) (*TokenPair, *domain.User, error) {
	device, err := s.verifyRedis.Verified(ctx, redis.VerifyEmailEntryKey, verifyToken)
	if err != nil {
		device, err = s.verifyRedis.Verified(ctx, redis.VerifySmsEntryKey, verifyToken)
		if err != nil {
			return nil, nil, fmt.Errorf("auth_service.LoginVerify: invalid or expired verify token: %w", err)
		}
	}

	var user *domain.User
	switch {
	case device.Email != nil && *device.Email != "":
		user, err = s.userRepo.GetByEmail(ctx, *device.Email)
	case device.PhoneNumber != nil && *device.PhoneNumber != "":
		user, err = s.userRepo.GetByPhoneNumber(ctx, *device.PhoneNumber)
	default:
		return nil, nil, fmt.Errorf("auth_service.LoginVerify: no verified contact method")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("auth_service.LoginVerify: failed to get user: %w", err)
	}

	tokens, err := s.createTokenPair(user.Username, user.UID, user.Role)
	if err != nil {
		return nil, nil, fmt.Errorf("auth_service.LoginVerify: failed to create tokens: %w", err)
	}

	clientIP, userAgent := getMetadataFromContext(ctx)
	if err := s.createSession(ctx, user.Username, tokens.RefreshToken, clientIP, userAgent, tokens.RefreshTokenPayload.ExpiredAt); err != nil {
		return nil, nil, fmt.Errorf("auth_service.LoginVerify: failed to create session: %w", err)
	}

	_ = s.verifyRedis.DelVerified(ctx, redis.VerifyEmailEntryKey, verifyToken)
	_ = s.verifyRedis.DelVerified(ctx, redis.VerifySmsEntryKey, verifyToken)

	if err := s.userRepo.UpdateLastLogin(ctx, user.Username); err != nil {
		return nil, nil, fmt.Errorf("auth_service.LoginVerify: failed to update last login: %w", err)
	}

	return tokens, user, nil
}

func (s *AuthService) Register(ctx context.Context, params domain.CreateUserParams) (*TokenPair, *domain.User, error) {
	clientIP, userAgent := getMetadataFromContext(ctx) // cooldown 5 minutes for one IP address

	if block, ttl, err := s.userRedis.AuthBlock(ctx, clientIP); block || err != nil {
		return nil, nil, fmt.Errorf("%w: %f", domain.ErrAuthBlock, ttl)
	}

	if params.Username == nil {
		v := generateNewUsername()
		existsByUsername, err := s.userRepo.ExistsByUsername(ctx, v)
		if err != nil {
			return nil, nil, fmt.Errorf("auth_service.Register: failed to check new username: %w", err)
		}
		for existsByUsername {
			v = generateNewUsername()
			existsByUsername, err = s.userRepo.ExistsByUsername(ctx, v)
			if err != nil {
				return nil, nil, fmt.Errorf("auth_service.Register: failed to check new username: %w", err)
			}
		}
		params.Username = &v
	}
	user, err := s.userRepo.Create(ctx, &params)
	if err != nil {
		return nil, nil, err
	}

	tokens, err := s.createTokenPair(user.Username, user.UID, "user")
	if err != nil {
		return nil, nil, fmt.Errorf("auth_service.Register: failed to create access and refresh tokens: %w", err)
	}

	if err := s.userRedis.SetProfile(ctx, user.Username, user); err != nil {
		return nil, nil, fmt.Errorf("auth_service.Register: failed to set profile in Redis: %w", err)
	}

	if err = s.createSession(ctx, user.Username, tokens.RefreshToken, clientIP, userAgent, tokens.RefreshTokenPayload.ExpiredAt); err != nil {
		return nil, nil, fmt.Errorf("auth_service.Register: failed to create session: %w", err)
	}

	if err := s.userRedis.SetAuthBlock(ctx, clientIP); err != nil {
		return nil, nil, fmt.Errorf("auth_service.Login: failed to set auth block in Redis: %w", err)
	}

	return tokens, user, nil
}

func (s *AuthService) RequestVerifyEmail(ctx context.Context, email string, operationId string) (string, error) {
	if exists, err := s.userRepo.ExistsByEmail(ctx, email); !exists {
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

	if err := s.verifyRedis.IncDailyByMethod(ctx, &domain.MethodToVerify{Email: &email}); err != nil {
		return "", fmt.Errorf("auth_service.RequestVerifyEmail: daily limit exceeded: %w", err)
	}

	if err := s.verifyRedis.SetCode(ctx, operationId, domain.NewVerifyState(
		&domain.MethodToVerify{Email: &email}, s.sender.Helper().Code2Hash(code), s.verifyRedis.Attempts(), s.verifyRedis.CodeTTL())); err != nil {
		return "", fmt.Errorf("auth_service.RequestVerifyEmail: failed to set verify code in Redis: %w", err)
	}

	if err := s.sender.SendEmail(email, code); err != nil {
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

	// Determine contact ID and entry code for cooldown/daily tracking
	var contactID string
	var entryCode string
	if method.PhoneNumber != nil && *method.PhoneNumber != "" {
		contactID = *method.PhoneNumber
		entryCode = redis.VerifySmsEntryKey
	} else if method.Email != nil && *method.Email != "" {
		contactID = *method.Email
		entryCode = redis.VerifyEmailEntryKey
	}

	// Check cooldown between attempts
	if contactID != "" {
		if cd := s.verifyRedis.CheckCooldown(ctx, entryCode, contactID); cd > 0 {
			return "", fmt.Errorf("%w: cooldown remaining %.0fs", domain.ErrSmsCooldown, cd.Seconds())
		}
	}

	if ok := s.sender.Helper().CompareCode(code, record.HashCode); !ok {
		// Decrement attempts and set cooldown on failure
		remaining, decErr := s.verifyRedis.DecrementAttempts(ctx, operationId)
		if decErr != nil {
			// All attempts exhausted
			return "", domain.ErrSmsRateLimited
		}
		if contactID != "" {
			_ = s.verifyRedis.SetCooldown(ctx, entryCode, contactID, s.verifyRedis.CooldownDuration())
		}
		if remaining <= 0 {
			return "", domain.ErrSmsRateLimited
		}
		return "", domain.ErrInvalidCode
	}

	_ = s.verifyRedis.DelCode(ctx, operationId)

	token := uuid.NewString()

	if method.Email != nil && *method.Email != "" {
		entryCode = redis.VerifyEmailEntryKey
	} else {
		entryCode = redis.VerifySmsEntryKey
	}

	if err := s.verifyRedis.SetVerified(ctx, entryCode, token, &method); err != nil {
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

func (s *AuthService) createSendCodeWebhook(ctx context.Context, phoneNumber, clientIP string) (string, error) {
	if err := s.verifyRedis.IncDailyByMethod(ctx, &domain.MethodToVerify{PhoneNumber: &phoneNumber}); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeWebhook: daily limit exceeded: %w", err)
	}
	if err := s.verifyRedis.IncDailyByIP(ctx, redis.VerifySmsEntryKey, clientIP); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeWebhook: daily IP limit exceeded: %w", err)
	}

	code := s.sender.Helper().GenerateVerifyCode()
	if err := s.sender.SendPhone(phoneNumber, code); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeWebhook: failed to send sms: %w", err)
	}

	operationId := redis.VerifySmsEntryKey + phoneNumber
	if err := s.verifyRedis.SetCode(
		ctx,
		operationId,
		domain.NewVerifyState(&domain.MethodToVerify{PhoneNumber: &phoneNumber}, s.sender.Helper().Code2Hash(code), s.verifyRedis.Attempts(), s.verifyRedis.CodeTTL()),
	); err != nil {
		return "", fmt.Errorf("auth_service.createSendCodeWebhook: failed to save sms verification code (redis): %w", err)
	}

	return operationId, nil
}

func (s *AuthService) Verify(ctx context.Context, verifyToken string) (*TokenPair, error) {
	device, err := s.verifyRedis.Verified(ctx, redis.VerifyEmailEntryKey, verifyToken)
	if err != nil {
		return nil, fmt.Errorf("auth_service.VerifyEmail: failed to verify email: %w", err)
	}

	user, err := s.userRepo.Verify(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("auth_service.Verify: failed to verify device: %w", err)
	}

	return s.createTokenPair(user.Username, user.UID, user.Role)
}

// createTokenPair
func (s *AuthService) createTokenPair(username, uid, role string) (*TokenPair, error) {
	accessToken, err := s.accessManager.CreateToken(username, uid, role)
	if err != nil {
		return nil, fmt.Errorf("auth_service.createTokenPair: create access token: %w", err)
	}

	refreshToken, payload, err := s.refreshManager.CreateRefreshToken(username, uid, role)
	if err != nil {
		return nil, fmt.Errorf("auth_service.createTolenPair: create refresh token: %w", err)
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, RefreshTokenPayload: payload}, nil
}

func (s *AuthService) createSession(ctx context.Context, username, refreshToken, clientIP, userAgent string, expiredAt time.Time) error {
	// Select session by username because only username store old data for old refreshToken
	// new refresh token while not store in database
	// it's only will be created after check old session
	session, err := s.userRedis.Session(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrRedisNotFound) {
			session, err = s.userRepo.GetSessionByUsername(ctx, username)
			if err != nil && !errors.Is(err, domain.ErrNotFound) {
				return fmt.Errorf("auth_service.createSession get old session: %w", err)
			}
		} else {
			return fmt.Errorf("auth_service.createSession get old session (redis): %w", err)
		}
	}

	// TODO: implement IPs and useragent checker
	if session != nil && session.ClientIP != clientIP && session.UserAgent != userAgent {
		return domain.ErrDeviceMistake
	}

	// Creating session record about user session in main database
	newSession, err := s.userRepo.CreateSession(ctx, domain.CreateSessionParams{
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    expiredAt,
	})
	if err != nil {
		return fmt.Errorf("auth_service.createSession create session: %w", err)
	}

	// Creating session record about user session in Redis
	if err := s.userRedis.SetSession(ctx, newSession); err != nil {
		return fmt.Errorf("auth_service.createSession create session (redis): %w", err)
	}

	return nil
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
