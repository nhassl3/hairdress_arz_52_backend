package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nhassl3/hairdress_arz/pkg/sms"
)

type User struct {
	Username    string    `json:"username"`
	UID         string    `json:"uid"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	IsVerified  bool      `json:"is_verified"`
	LastLogin   time.Time `json:"last_login"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	if u == nil {
		return ErrRedisNotFound
	}
	return json.Unmarshal(data, u)
}

type SmsCodeRecorder struct {
	Hash         string    `json:"hash"` // hex string SHA-256-HMAC
	AttemptsLeft int32     `json:"attempts_left"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewSmsRecord(hash string, attemptsLeft int32, ttl time.Duration) *SmsCodeRecorder {
	return &SmsCodeRecorder{
		Hash:         hash,
		AttemptsLeft: attemptsLeft,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(ttl),
	}
}

func (s *SmsCodeRecorder) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SmsCodeRecorder) UnmarshalBinary(data []byte) error {
	if s == nil {
		return ErrRedisNotFound
	}
	return json.Unmarshal(data, s)
}

type CreateUserParams struct {
	Username    *string
	FullName    *string
	PhoneNumber string
}

type UserRepository interface {
	Create(ctx context.Context, params *CreateUserParams) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
	Verify(ctx context.Context, username string) error
	UpdateLastLogin(ctx context.Context, username string) error
}

// Redis interfaces

type UserRedis interface {
	Profile(ctx context.Context, username string) (*User, error)
	SetProfile(ctx context.Context, username string, profile *User) error
	DelProfile(ctx context.Context, username string) error
	AuthBlock(ctx context.Context, clientIP string) (bool, float64, error)
	SetAuthBlock(ctx context.Context, clientIP string) error
	DelAuthBlock(ctx context.Context, clientIP string) error
}

type SmsRedis interface {
	SaveCode(ctx context.Context, phone, hash string) (*SmsCodeRecorder, error)
	GetCode(ctx context.Context, phone string) (*SmsCodeRecorder, error)
	DeleteCode(ctx context.Context, phone string) error
	DecrementAttempts(ctx context.Context, phone string) (left int32, err error)
	CheckCooldown(ctx context.Context, phone string) time.Duration
	SetCooldown(ctx context.Context, phone string, duration time.Duration) error
	IncDailyByPhoneNumber(ctx context.Context, phoneNumber string) error
	IncDailyByIP(ctx context.Context, ip string) error
}

type SmsSender interface {
	Send(ctx context.Context, phone, code string) error
	Helper() *sms.Helper
}
