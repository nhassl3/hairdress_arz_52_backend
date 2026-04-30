package service

import (
	"context"
	"fmt"
	"strings"

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
}

type TokenPair struct {
	AccessToken,
	RefreshToken string
	RefreshTokenPayload *authServiceHub.Payload
}

func NewAuthService(userRepo domain.UserRepository, accessManager, refreshManager auth.TokenManager,
	blacklist auth.TokenBlacklist, userRedis domain.UserRedis) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		refreshManager: refreshManager,
		accessManager:  accessManager,
		blacklist:      blacklist,
		userRedis:      userRedis,
	}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	return "", nil
}

func (s *AuthService) Register(ctx context.Context, params *domain.CreateUserParams) (*domain.User, *TokenPair, error) {
	// cooldown 5 minutes for one IP address
	clientIP, _ := getMetadataFromContext(ctx)

	if block, ttl, err := s.userRedis.AuthBlock(ctx, clientIP); block || err != nil {
		return nil, nil, fmt.Errorf("%w: %f", domain.ErrAuthBlock, ttl)
	}

	if params.Username == nil {
		v := generateNewUsername()
		existsByUsername, err := s.userRepo.ExistsByUsername(ctx, v)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to check new username: %w", err)
		}
		for existsByUsername {
			v = generateNewUsername()
			existsByUsername, err = s.userRepo.ExistsByUsername(ctx, v)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to check new username: %w", err)
			}
		}
		params.Username = &v
	}
	user, err := s.userRepo.Create(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	// TODO: sms code business logic

	// TODO: create pair of tokens
	//tokens, err := s.createTokenPair(user.Username, uuid.New().String())

	return user, nil, nil
}

func (s *AuthService) Logout(ctx context.Context) {
}

func (s *AuthService) RefreshToken(ctx context.Context) string {
	return ""
}

func (s *AuthService) GetMe(ctx context.Context) (*domain.User, error) {
	return nil, nil
}

// helpers

func generateNewUsername() string {
	return fmt.Sprintf("@user_%d", uuid.New().ID())
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
