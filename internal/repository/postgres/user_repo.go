package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhassl3/hairdress_arz/internal/db"
	"github.com/nhassl3/hairdress_arz/internal/domain"
)

type AuthRepo struct {
	store *db.Store
}

func NewAuthRepository(store *db.Store) *AuthRepo {
	return &AuthRepo{
		store: store,
	}
}

// User repository functions

func (repo *AuthRepo) Create(ctx context.Context, params *domain.CreateUserParams) (*domain.User, error) {
	var user db.User

	if err := repo.store.ExecTx(ctx, func(q *db.Queries) error {
		var fnErr error

		existsByPhone, fnErr := q.ExistsByPhoneNumber(ctx, params.PhoneNumber)
		if fnErr != nil {
			return fmt.Errorf("failed to find user by phone number: %w", fnErr)
		}
		if existsByPhone {
			return domain.ErrPhoneAlreadyExists
		}

		// additional check
		existsByUsername, fnErr := q.ExistsByUsername(ctx, *params.Username)
		if fnErr != nil {
			return fmt.Errorf("failed to find user by username: %w", fnErr)
		}
		if existsByUsername {
			return domain.ErrUsernameAlreadyExists
		}

		user, fnErr = q.CreateUser(ctx, db.CreateUserParams{
			Username:    str2Text(params.Username),
			Email:       params.Email,
			PhoneNumber: params.PhoneNumber,
		})
		if fnErr != nil {
			// this check needed for that if account created twice
			// and other attempt just not provide
			// it's ok if in frontend may be a few requests in one method
			if pgErr, ok := errors.AsType[*pgconn.PgError](fnErr); ok {
				if pgErr.Code == "23505" { // unique constraint
					return domain.ErrUsernameAlreadyExists
				}
			}
			return fnErr
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return toDomainUser(&user), nil
}

func (repo *AuthRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := repo.store.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by phone number: %w", err)
	}
	return toDomainUser(&user), nil
}

func (repo *AuthRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := repo.store.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("user_repo.GetByEmail: failed to get user by email: %w", err)
	}
	return toDomainUser(&user), nil
}

func (repo *AuthRepo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	user, err := repo.store.GetUserByPhone(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by phone number: %w", err)
	}
	return toDomainUser(&user), nil
}

func (repo *AuthRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return repo.store.ExistsByUsername(ctx, username)
}

func (repo *AuthRepo) ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error) {
	return repo.store.ExistsByPhoneNumber(ctx, phoneNumber)
}

func (repo *AuthRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return repo.store.ExistsByEmail(ctx, email)
}

func (repo *AuthRepo) Verify(ctx context.Context, toVerify *domain.MethodToVerify) (*domain.User, error) {
	user, err := repo.store.VerifyUser(ctx, db.VerifyUserParams{
		PhoneNumber: str2Text(toVerify.PhoneNumber),
		Email:       str2Text(toVerify.Email),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("user_repo.Verify: failed to verify user: %w", err)
	}
	return toDomainUser(&user), nil
}

func (repo *AuthRepo) UpdateLastLogin(ctx context.Context, username string) error {
	return repo.store.UpdateLastLogin(ctx, username)
}

// Session repository functions

func (repo *AuthRepo) CreateSession(ctx context.Context, params domain.CreateSessionParams) (*domain.Session, error) {
	session, err := repo.store.CreateSession(ctx, db.CreateSessionParams{
		Username:     params.Username,
		RefreshToken: params.RefreshToken,
		UserAgent:    params.UserAgent,
		ClientIp:     params.ClientIp,
		ExpiresAt:    pgtype.Timestamptz{Time: params.ExpiresAt, Valid: true},
		IsBlocked:    params.IsBlocked,
	})
	if err != nil {
		return nil, fmt.Errorf("user_repo.CreateSession: failed to create session: %w", err)
	}
	return toDomainSession(session), nil
}

func (repo *AuthRepo) GetSession(ctx context.Context, refreshToken string) (*domain.Session, error) {
	row, err := repo.store.GetSession(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("user_repo.GetSession: %w", err)
	}
	return toDomainSession(row), nil
}

func (repo *AuthRepo) GetSessionByUsername(ctx context.Context, username string) (*domain.Session, error) {
	row, err := repo.store.GetSessionByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("user_repo.GetSessionByUsername: %w", err)
	}
	return toDomainSession(row), nil
}

func (repo *AuthRepo) DeleteSession(ctx context.Context, username string) error {
	if err := repo.store.DeleteSession(ctx, username); err != nil {
		return fmt.Errorf("user_repo.DeleteSession: %w", err)
	}
	return nil
}
