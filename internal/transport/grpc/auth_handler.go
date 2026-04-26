package grpc

import (
	"context"
	"time"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/nhassl3/hairdress_arz/internal/service"
	authv1 "github.com/nhassl3/hairdress_arz_52_contracts/pkg/pb/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {

	return &authv1.LoginResponse{}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	return nil, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	return nil, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	return nil, nil
}

func (h *AuthHandler) GetMe(ctx context.Context, req *authv1.GetMeRequest) (*authv1.GetMeResponse, error) {
	return nil, nil
}

// proto

func protoUser(user *domain.User) *authv1.UserInfo {
	return &authv1.UserInfo{
		Username:    user.Username,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		IsVerified:  user.IsVerified,
		CreatedAt:   safeTimestamp(user.CreatedAt),
		UpdatedAt:   safeTimestamp(user.UpdatedAt),
	}
}

// safeTimestamp converts a time.Time to a protobuf Timestamp.
// Zero times are returned as nil.
func safeTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}
