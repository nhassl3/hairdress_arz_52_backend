package domain

import (
	"context"
	"encoding/json"
	"time"
)

type User struct {
	Username    string    `json:"username"`
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

	// VerifyAndTouch two operation in function
	VerifyAndTouch(ctx context.Context, username string) error
}

type UserRedis interface {
	Profile(ctx context.Context, username string) (*User, error)
	SetProfile(ctx context.Context, username string, profile *User) error
	DelProfile(ctx context.Context, username string) error
	AuthBlock(ctx context.Context, clientIP string) (bool, float64, error)
	SetAuthBlock(ctx context.Context, clientIP string) error
	DelAuthBlock(ctx context.Context, clientIP string) error
}

type SmsRepository interface {
	SaveCode(ctx context.Context, phone, hash string, attempts int32, ttl time.Duration) (*SmsCodeRecorder, error)
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
}
