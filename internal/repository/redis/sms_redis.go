package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/redis/go-redis/v9"
)

// Keys for SMS Redis
const (
	CodePrefix     = "sms:code:"        // <phone>
	AttemptsPrefix = "sms:attempts:"    // <phone>
	CooldownPrefix = "sms:cooldown:"    // <phone> - empty value, only TTL
	PhoneDaily     = "sms:daily:phone:" // <phone>:<YYYY-MM-DD>
	IPDaily        = "sms:daily:ip:"    // <IP>:<YYYY-MM-DD>
)

// LUA scripts
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
if ok then return 0 end
return decoded.attempts_left
`

const incByPhoneScript = `
local v = redis.call("GET", KEYS[1])
if not v then
	redis.call("SET", KEYS[1], 1)
	redis.call("EXPIRE", KEYS[1], 86400)
	return redis.status_reply('OK')
end
if v >= ARGV[1] then
	return redis.error_reply("out of daily limits for phone number")
end
redis.call("SET", KEYS[1], v+1)
return redis.status_reply('OK')
`

const incByIPScript = `
local v = redis.call("GET", KEYS[1])
if not v then
	redis.call("SET", KEYS[1], ARGV[1])
	redis.call("EXPIRE", KEYS[1], 86400)
	return redis.status_reply('OK')
end
if v >= ARGV[1] then
	return redis.error_reply("out of daily limits for ip address")
end
redis.call("SET", KEYS[1], v+1)
return redis.status_reply('OK')
`

type SMSRedis struct {
	client *redis.Client
	SMSVerificationCodeTTL,
	SMSCooldown time.Duration
	Attempts,
	PhoneDailyAttempts,
	IPDailyAttempts int32
}

func NewSMSRedis(client *redis.Client, smsVerificationCodeTTL, smsCooldown time.Duration, attempts, phoneDailyAttempts, ipDailyAttempts int32) *SMSRedis {
	return &SMSRedis{
		client:                 client,
		SMSVerificationCodeTTL: smsVerificationCodeTTL,
		SMSCooldown:            smsCooldown,
		Attempts:               attempts,
		PhoneDailyAttempts:     phoneDailyAttempts,
		IPDailyAttempts:        ipDailyAttempts,
	}
}

// For six-digit code prefix

// SaveCode create a record in Redis storage with six-digit code for given phone number
func (r *SMSRedis) SaveCode(ctx context.Context, phone, hash string) (*domain.SmsCodeRecorder, error) {
	smsRecord := domain.NewSmsRecord(hash, r.Attempts, r.SMSVerificationCodeTTL)
	if err := r.client.SetArgs(ctx,
		CodePrefix+phone,
		smsRecord,
		redis.SetArgs{
			Mode: "NX",
			TTL:  r.SMSVerificationCodeTTL,
		},
	).Err(); err != nil {
		return nil, fmt.Errorf("failed to save code for phone %s: %w", phone, err)
	}
	return smsRecord, nil
}

// GetCode returns struct with hash, attempts left, createdAt and TTL for given phone number
func (r *SMSRedis) GetCode(ctx context.Context, phone string) (*domain.SmsCodeRecorder, error) {
	var smsRecord domain.SmsCodeRecorder
	if err := r.client.Get(ctx, CodePrefix+phone).Scan(&smsRecord); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrRedisNotFound
		}
		return nil, fmt.Errorf("failed to get code for phone %s: %w", phone, err)
	}
	// ordinary Redis clean key:value after given ttl
	if time.Now().After(smsRecord.ExpiresAt) {
		return nil, domain.ErrRedisCodeExpired
	}
	return &smsRecord, nil
}

// DeleteCode remove record from Redis database with code for phone number
func (r *SMSRedis) DeleteCode(ctx context.Context, phone string) error {
	return r.client.Del(ctx, CodePrefix+phone).Err()
}

// methods with struct of code prefix

func (r *SMSRedis) DecrementAttempts(ctx context.Context, phone string) (left int32, err error) {
	v, err := r.client.Eval(ctx, decrementScript, []string{CodePrefix + phone}, nil).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement attempts for phone %s: %w", phone, err)
	}
	if v == 0 {
		return 0, domain.ErrRedisNotFound
	}
	return left, nil
}

// cooldown

func (r *SMSRedis) CheckCooldown(ctx context.Context, phone string) time.Duration {
	return r.client.TTL(ctx, CooldownPrefix+phone).Val()
}

func (r *SMSRedis) SetCooldown(ctx context.Context, phone string, duration time.Duration) error {
	return r.client.Set(ctx, CooldownPrefix+phone, nil, duration).Err()
}

// Increment daily attempts for user by phone number

func (r *SMSRedis) IncDailyByPhoneNumber(ctx context.Context, phoneNumber string) error {
	return r.client.Eval(ctx, incByPhoneScript, []string{PhoneDaily + phoneNumber}, r.PhoneDailyAttempts).Err()
}

// Increment daily attempts for user by IP address

func (r *SMSRedis) IncDailyByIP(ctx context.Context, ip string) error {
	return r.client.Eval(ctx, incByIPScript, []string{IPDaily + ip}, r.IPDailyAttempts).Err()
}
