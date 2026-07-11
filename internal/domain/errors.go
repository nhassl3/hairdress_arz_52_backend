package domain

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrPhoneAlreadyExists    = errors.New("phone already exists")
	ErrForbidden             = errors.New("forbidden")
	ErrInvalidCode           = errors.New("invalid code")
	ErrCodeExpired           = errors.New("code expired")
	ErrTooManyAttempts       = errors.New("too many attempts")
	ErrSmsRateLimited        = errors.New("sms rate limited")
	ErrSessionIsBlocked      = errors.New("session is blocked")
	ErrSmsCooldown           = errors.New("sms cooldown")
	ErrInvalidToken          = errors.New("invalid token")
	ErrUserNotVerified       = errors.New("user not verified")
	ErrRedisNotFound         = errors.New("redis not found")
	ErrAuthBlock             = errors.New("auth temporally blocked")
	ErrRedisCodeExpired      = errors.New("sms code expired")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrNotFound              = errors.New("record not found")
	ErrExpiredToken          = errors.New("refresh token is expired")
	ErrDeviceMistake         = errors.New("device mistake")
	ErrInvalidRequestMethod  = errors.New("invalid request method")
	ErrNoBookings            = errors.New("no bookings")
	ErrDataNoProvide         = errors.New("data not provide")
)
