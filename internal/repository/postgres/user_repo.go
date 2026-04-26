package postgres

import "github.com/nhassl3/hairdress_arz/internal/db"

type AuthRepo struct {
	store *db.Store
}

Login(ctx context.Context, username string, password string) (User, error)
Register(ctx context.Context, username string, password string) (User, error)
Logout(ctx context.Context, username string) error
RefreshToken(ctx context.Context, username string) (string, error)
GetMe(ctx context.Context) (*User, error)
