package grpc

import (
	"context"
	"time"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/nhassl3/hairdress_arz/internal/service"
	"github.com/nhassl3/hairdress_arz/internal/transport/grpc/interceptors"
	authv1 "github.com/nhassl3/hairdress_arz_52_contracts/pkg/pb/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	user, err := h.svc.Login(ctx, req.GetPhoneNumber())
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.LoginResponse{
		User: protoUser(user),
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, err := h.svc.Register(ctx, domain.CreateUserParams{
		Username:    &req.Username,
		FullName:    nil,
		PhoneNumber: req.GetPhoneNumber(),
	})
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.RegisterResponse{
		User: protoUser(user),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, _ *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	payload, ok := interceptors.PayloadFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "missing token")
	}

	if err := h.svc.Logout(ctx, payload); err != nil {
		return nil, domainErr(err)
	}

	return &authv1.LogoutResponse{Success: true}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	tokens, err := h.svc.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, domainErr(err)
	}
	if tokens == nil {
		return &authv1.RefreshTokenResponse{}, nil
	}
	return &authv1.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (h *AuthHandler) GetMe(ctx context.Context, req *authv1.GetMeRequest) (*authv1.GetMeResponse, error) {
	username := interceptors.GetNameFromContext(ctx)
	if username == "" {
		return nil, status.Error(codes.FailedPrecondition, "missing token")
	}
	user, err := h.svc.GetMe(ctx, username)
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.GetMeResponse{
		User: protoUser(user),
	}, nil
}

func (h *AuthHandler) VerifyCode(ctx context.Context, req *authv1.VerifyCodeRequest) (*authv1.VerifyCodeResponse, error) {
	tokenPair, err := h.svc.VerifyCode(ctx, req.GetPhoneNumber(), req.GetCode())
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.VerifyCodeResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
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
