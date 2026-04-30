package redis

import (
	"context"
	"errors"
	"time"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	profilePrefix   = "profile:"
	authBlockPrefix = "auth:block:"
)

type UserRedis struct {
	client       *redis.Client
	profileTTL   time.Duration
	authBlockTTL time.Duration
}

func NewUserRedis(client *redis.Client, profileTTL, authBlockTTL time.Duration) *UserRedis {
	return &UserRedis{
		client:       client,
		profileTTL:   profileTTL,
		authBlockTTL: authBlockTTL,
	}
}

func (u *UserRedis) Profile(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	if err := u.client.Get(ctx, profilePrefix+username).Scan(&user); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrRedisNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRedis) SetProfile(ctx context.Context, username string, profile *domain.User) error {
	return u.client.Set(ctx, profilePrefix+username, &profile, u.profileTTL).Err()
}

func (u *UserRedis) DelProfile(ctx context.Context, username string) error {
	return u.client.Del(ctx, profilePrefix+username).Err()
}

func (u *UserRedis) AuthBlock(ctx context.Context, clientIP string) (bool, float64, error) {
	ok, err := u.client.Get(ctx, authBlockPrefix+clientIP).Bool()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, -1, nil // -1 - not records
		}
		return false, -1, err // not records but some other errors catches
	}
	ttl := u.client.TTL(ctx, authBlockPrefix+clientIP).Val()
	return ok, ttl.Seconds(), nil
}

func (u *UserRedis) SetAuthBlock(ctx context.Context, clientIP string) error {
	return u.client.Set(ctx, authBlockPrefix+clientIP, true, u.authBlockTTL).Err()
}

func (u *UserRedis) DelAuthBlock(ctx context.Context, clientIP string) error {
	return u.client.Del(ctx, authBlockPrefix+clientIP).Err()
}
