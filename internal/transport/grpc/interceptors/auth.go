package interceptors

import (
	"context"
	"slices"
	"strings"

	"github.com/nhassl3/servicehub-backend/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const PayloadKey contextKey = "auth_payload"

// publicMethods lists gRPC methods that do not require authentication
var publicMethods = []string{
	"/auth.v1.AuthService/Register",
	"/auth.v1.AuthService/Login",
	"/auth.v1.AuthService/Logout",
	"/auth.v1.AuthService/RefreshToken",
	"/auth.v1.AuthService/ApproveCode",
	"/auth.v1.AuthService/LoginVerify",
	"/booking.v1.BookingService/CreateBooking",
	"/booking.v1.BookingService/GetBooking",
}

// AuthInterceptor returns a gRPC unary interceptor for PASETO token verification
func AuthInterceptor(tokenManager auth.TokenManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if slices.Contains(publicMethods, info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization is not provided")
		}

		authHeader := values[0]
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := tokenManager.VerifyToken(token)
		if err != nil {
			if auth.IsAny(err) {
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}

		ctx = context.WithValue(ctx, PayloadKey, payload)
		return handler(ctx, req)
	}
}

// PayloadFromContext extracts the auth payload from context
func PayloadFromContext(ctx context.Context) (*auth.Payload, bool) {
	payload, ok := ctx.Value(PayloadKey).(*auth.Payload)
	return payload, ok
}

func GetNameFromContext(ctx context.Context) string {
	payload, ok := ctx.Value(PayloadKey).(*auth.Payload)
	if !ok {
		return ""
	}
	return payload.Username
}
