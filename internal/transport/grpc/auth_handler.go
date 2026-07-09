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

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthResponse, error) {
	var loginParams *domain.LoginParams
	switch v := req.GetMethod().(type) {
	case *authv1.LoginRequest_Email:
		loginParams = &domain.LoginParams{Email: &v.Email}
	case *authv1.LoginRequest_Username:
		loginParams = &domain.LoginParams{Username: &v.Username}
	case *authv1.LoginRequest_PhoneNumber:
		loginParams = &domain.LoginParams{PhoneNumber: &v.PhoneNumber}
	default:
		return nil, domainErr(domain.ErrInvalidRequestMethod)
	}

	operationId, user, err := h.svc.Login(ctx, loginParams)
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.AuthResponse{
		User:        protoUser(user),
		OperationId: &operationId,
	}, nil
}

func (h *AuthHandler) LoginVerify(ctx context.Context, req *authv1.LoginVerifyRequest) (*authv1.LoginVerifyResponse, error) {
	tokens, user, err := h.svc.LoginVerify(ctx, req.GetVerifyToken())
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.LoginVerifyResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         protoUser(user),
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AuthResponse, error) {
	tokens, user, err := h.svc.Register(ctx, domain.CreateUserParams{
		Username:    &req.Username,
		Email:       req.GetEmail(),
		PhoneNumber: req.GetPhoneNumber(),
	})
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         protoUser(user),
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

func (h *AuthHandler) ApproveCode(ctx context.Context, req *authv1.ApproveCodeRequest) (*authv1.ApproveCodeResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var method domain.MethodToVerify
	switch v := req.GetMethod().(type) {
	case *authv1.ApproveCodeRequest_Email:
		method.Email = &v.Email
	case *authv1.ApproveCodeRequest_PhoneNumber:
		method.PhoneNumber = &v.PhoneNumber
	default:
		return nil, domainErr(domain.ErrInvalidRequestMethod)
	}

	token, err := h.svc.ApproveCode(ctx, req.GetOperationId(), method, req.GetCode())
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.ApproveCodeResponse{
		Token: token,
	}, nil
}

func (h *AuthHandler) RequestEmailVerify(ctx context.Context, req *authv1.RequestEmailVerifyRequest) (*authv1.RequestEmailVerifyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	operationId, err := h.svc.RequestVerifyEmail(ctx, req.GetEmail(), req.GetOperationId())
	if err != nil {
		return nil, domainErr(err)
	}

	return &authv1.RequestEmailVerifyResponse{
		OperationId: operationId,
	}, nil
}

func (h *AuthHandler) VerifyEmail(ctx context.Context, req *authv1.VerifyEmailRequest) (*authv1.VerifyEmailResponse, error) {
	tokenPair, err := h.svc.Verify(ctx, req.GetVerifyToken())
	if err != nil {
		return nil, domainErr(err)
	}
	return &authv1.VerifyEmailResponse{
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
		Email:       user.Email,
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
