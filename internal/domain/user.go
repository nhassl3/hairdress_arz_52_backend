package domain

import (
	"context"
	"time"
)

type User struct {
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserRepository interface {
	Login(ctx context.Context, username string, password string) (User, error)
	Register(ctx context.Context, username string, password string) (User, error)
	Logout(ctx context.Context, username string) error
	RefreshToken(ctx context.Context, username string) (string, error)
	GetMe(ctx context.Context) (*User, error)
}

type UserRedis interface {
	Profile(ctx context.Context, username string) (User, error)
	SetProfile(ctx context.Context, username string, profile User) error
}
