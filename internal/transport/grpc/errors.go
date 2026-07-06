package grpc

import (
	"errors"
	"fmt"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// domainErr maps a domain error to the appropriate gRPC status error.
// Unknown errors are returned as Internal.
func domainErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrUsernameAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrPhoneAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrCodeExpired):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, domain.ErrUserNotVerified):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, domain.ErrInvalidCode):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrTooManyAttempts):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrSmsRateLimited):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrSmsCooldown):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrRedisNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidToken):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrAuthBlock):
		return status.Error(codes.Unavailable, err.Error())
	default:
		return status.Error(codes.Internal, fmt.Sprintf("internal server error: %s", err.Error()))
	}
}
