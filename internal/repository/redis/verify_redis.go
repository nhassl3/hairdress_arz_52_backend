package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	codePrefix     = "code:"
	CooldownPrefix = "code:cooldown:"
	dailyPrefix    = "code:daily:"
	IPDaily        = "code:daily:ip:"

	verifiedPrefix = "verified:"

	VerifyEmailEntryKey = "verify_email:"
	VerifySmsEntryKey   = "verify_sms:"
)

const decrementScript = `
local v = redis.call("GET", KEYS[1])
if not v then return 0 end
local ok, decoded = pcall(cjson.decode, v)
if not ok then return 0 end
if decoded.attempts_left <= 0 then
	return redis.call("DEL", KEYS[1])
end
decoded.attempts_left--
local ok, new_json_str = pcall(cjson.encode, decoded)
if not ok then return 0 end
return decoded.attempts_left
`

const incDailyScript = `
local v = redis.call("GET", KEYS[1])
if not v then
	redis.call("SET", KEYS[1], 1)
	redis.call("EXPIRE", KEYS[1], 86400)
	return redis.status_reply('OK')
end
if v >= ARGV[1] then
	return redis.error_reply("out of daily limits")
end
redis.call("SET", KEYS[1], v+1)
return redis.status_reply('OK')
`

const incByIPScript = `
local v = redis.call("GET", KEYS[1])
if not v then
	redis.call("SET", KEYS[1], 1)
	redis.call("EXPIRE", KEYS[1], 86400)
	return redis.status_reply('OK')
end
if v >= ARGV[1] then
	return redis.error_reply("out of daily limits for ip address")
end
redis.call("SET", KEYS[1], v+1)
return redis.status_reply('OK')
`

type VerifyRedis struct {
	client *redis.Client
	codeTTL,
	cooldown time.Duration
	attempts,
	phoneDailyAttempts,
	emailDailyAttempts,
	ipDailyAttempts int32
}

func NewVerifyRedis(client *redis.Client, codeTTL, cooldown time.Duration, attempts, phoneDailyAttempts, emailDailyAttempts, ipDailyAttempts int32) *VerifyRedis {
	return &VerifyRedis{
		client:             client,
		codeTTL:            codeTTL,
		cooldown:           cooldown,
		attempts:           attempts,
		phoneDailyAttempts: phoneDailyAttempts,
		emailDailyAttempts: emailDailyAttempts,
		ipDailyAttempts:    ipDailyAttempts,
	}
}

// ── Code methods (unified for both SMS and email) ────────────────────────────

func (r *VerifyRedis) SetCode(ctx context.Context, operationId string, state *domain.VerifyState) error {
	return r.client.SetArgs(ctx,
		codePrefix+operationId,
		state,
		redis.SetArgs{
			Mode: "NX",
			TTL:  r.codeTTL,
		},
	).Err()
}

func (r *VerifyRedis) Code(ctx context.Context, operationId string) (*domain.VerifyState, error) {
	var state domain.VerifyState
	if err := r.client.Get(ctx, codePrefix+operationId).Scan(&state); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrRedisNotFound
		}
		return nil, fmt.Errorf("failed to get code for %s: %w", operationId, err)
	}
	if time.Now().After(state.ExpiresAt) {
		return nil, domain.ErrRedisCodeExpired
	}
	return &state, nil
}

func (r *VerifyRedis) DelCode(ctx context.Context, operationId string) error {
	return r.client.Del(ctx, codePrefix+operationId).Err()
}

// ── Attempts / cooldown (unified for both SMS and email) ─────────────────────

func (r *VerifyRedis) DecrementAttempts(ctx context.Context, operationId string) (int32, error) {
	v, err := r.client.Eval(ctx, decrementScript, []string{codePrefix + operationId}, nil).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement attempts for %s: %w", operationId, err)
	}
	attempts, ok := v.(int64)
	if !ok || attempts == 0 {
		return 0, domain.ErrRedisNotFound
	}
	return int32(attempts), nil
}

func (r *VerifyRedis) CodeTTL() time.Duration          { return r.codeTTL }
func (r *VerifyRedis) Attempts() int32                 { return r.attempts }
func (r *VerifyRedis) CooldownDuration() time.Duration { return r.cooldown }

func (r *VerifyRedis) CheckCooldown(ctx context.Context, entryCode, id string) time.Duration {
	return r.client.TTL(ctx, CooldownPrefix+entryCode+id).Val()
}

func (r *VerifyRedis) SetCooldown(ctx context.Context, entryCode, id string, duration time.Duration) error {
	return r.client.Set(ctx, CooldownPrefix+entryCode+id, nil, duration).Err()
}

// ── Daily rate limiting ──────────────────────────────────────────────────────

func (r *VerifyRedis) IncDailyByMethod(ctx context.Context, method *domain.MethodToVerify) error {
	if method.PhoneNumber != nil && *method.PhoneNumber != "" {
		return r.client.Eval(ctx, incDailyScript, []string{dailyPrefix + VerifySmsEntryKey + *method.PhoneNumber}, r.phoneDailyAttempts).Err()
	}
	if method.Email != nil && *method.Email != "" {
		return r.client.Eval(ctx, incDailyScript, []string{dailyPrefix + VerifyEmailEntryKey + *method.Email}, r.emailDailyAttempts).Err()
	}
	return nil
}

func (r *VerifyRedis) IncDailyByIP(ctx context.Context, entryCode, ip string) error {
	return r.client.Eval(ctx, incByIPScript, []string{IPDaily + entryCode + ip}, r.ipDailyAttempts).Err()
}

// ── Verified methods ─────────────────────────────────────────────────────────

func (r *VerifyRedis) Verified(ctx context.Context, entryCode, token string) (*domain.MethodToVerify, error) {
	var method domain.MethodToVerify
	if err := r.client.Get(ctx, verifiedPrefix+entryCode+token).Scan(&method); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrRedisNotFound
		}
		return nil, err
	}
	return &method, nil
}

func (r *VerifyRedis) SetVerified(ctx context.Context, entryCode, token string, method *domain.MethodToVerify) error {
	return r.client.Set(ctx, verifiedPrefix+entryCode+token, method, r.codeTTL).Err()
}

func (r *VerifyRedis) DelVerified(ctx context.Context, entryCode, token string) error {
	return r.client.Del(ctx, verifiedPrefix+entryCode+token).Err()
}
