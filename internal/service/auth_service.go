package service

import (
	"context"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/nhassl3/servicehub-backend/pkg/auth"
)

type AuthService struct {
	userRepo       domain.UserRepository
	userRedis      domain.UserRedis
	tokenManager   auth.TokenManager
	refreshManager auth.TokenManager
	blacklist      auth.TokenBlacklist
}

func NewAuthService(userRepo domain.UserRepository, userRedis domain.UserRedis,
	tokenManager, refreshManager auth.TokenManager,
	blacklist auth.TokenBlacklist) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		userRedis:      userRedis,
		tokenManager:   tokenManager,
		refreshManager: refreshManager,
		blacklist:      blacklist,
	}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (string, error) {
	return "", nil
}

func (s *AuthService) Register(ctx context.Context, user *domain.User) error {
	return nil
}

func (s *AuthService) Logout(ctx context.Context) {
}

func (s *AuthService) RefreshToken(ctx context.Context) string {
	return ""
}

func (s *AuthService) GetMe(ctx context.Context) (*domain.User, error) {
	return nil, nil
}
