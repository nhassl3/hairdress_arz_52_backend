package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nhassl3/hairdress_arz/internal/db"
	"github.com/nhassl3/hairdress_arz/internal/domain"
)

type AuthRepo struct {
	store *db.Store
}

func NewAuthRepo(store *db.Store) *AuthRepo {
	return &AuthRepo{
		store: store,
	}
}

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
			Username:    *params.Username,
			FullName:    str2Text(params.FullName),
			PhoneNumber: params.PhoneNumber,
		})
		if fnErr != nil {
			// this check needed for that if account created twice
			// and other attempt just not provide
			// it's ok if in frontend may be a few requests in one method
			var pgErr *pgconn.PgError
			if errors.As(fnErr, &pgErr) {
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

func (repo *AuthRepo) Verify(ctx context.Context, username string) error {
	return repo.store.VerifyUser(ctx, username)
}

func (repo *AuthRepo) UpdateLastLogin(ctx context.Context, username string) error {
	return repo.store.UpdateLastLogin(ctx, username)
}
