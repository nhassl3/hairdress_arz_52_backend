package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nhassl3/hairdress_arz/pkg/verify"
)

type User struct {
	Username    string    `json:"username"`
	UID         string    `json:"uid"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
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

type CreateUserParams struct {
	Username    *string
	FullName    *string
	PhoneNumber string
}

type MethodToVerify struct {
	PhoneNumber,
	Email *string
}

func (m *MethodToVerify) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MethodToVerify) UnmarshalBinary(data []byte) error {
	if m == nil {
		return ErrRedisNotFound
	}
	return json.Unmarshal(data, m)
}

type VerifyState struct {
	*MethodToVerify
	HashCode     string
	AttemptsLeft int32     `json:"attempts_left"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (r *VerifyState) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}

func (r *VerifyState) UnmarshalBinary(data []byte) error {
	if r == nil {
		return ErrRedisNotFound
	}
	return json.Unmarshal(data, r)
}

func NewVerifyState(method *MethodToVerify, hash string, attemptsLeft int32, ttl time.Duration) *VerifyState {
	return &VerifyState{
		MethodToVerify: method,
		HashCode:       hash,
		AttemptsLeft:   attemptsLeft,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(ttl),
	}
}

type UserRepository interface {
	Create(ctx context.Context, params *CreateUserParams) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhoneNumber(ctx context.Context, phoneNumber string) (bool, error)
	Verify(ctx context.Context, username string) error
	UpdateLastLogin(ctx context.Context, username string) error
	CreateSession(ctx context.Context, params CreateSessionParams) (*Session, error)
	GetSession(ctx context.Context, refreshToken string) (*Session, error)
	GetSessionByUsername(ctx context.Context, username string) (*Session, error)
	DeleteSession(ctx context.Context, username string) error
}

// Redis interfaces

type UserRedis interface {
	Profile(ctx context.Context, username string) (*User, error)
	SetProfile(ctx context.Context, username string, profile *User) error
	DelProfile(ctx context.Context, username string) error
	AuthBlock(ctx context.Context, clientIP string) (bool, float64, error)
	SetAuthBlock(ctx context.Context, clientIP string) error
	DelAuthBlock(ctx context.Context, clientIP string) error
	Session(ctx context.Context, username string) (*Session, error)
	SetSession(ctx context.Context, session *Session) error
	DelSession(ctx context.Context, username string) error
}

type VerifySender interface {
	SendPhone(phone, code string) error
	SendEmail(email, code string) error
	Helper() *verify.Helper
}

type VerifyRedis interface {
	Code(ctx context.Context, operationId string) (*VerifyState, error)
	SetCode(ctx context.Context, operationId string, code *VerifyState) error
	DelCode(ctx context.Context, operationId string) error
	Verified(ctx context.Context, entryCode, token string) (*MethodToVerify, error)
	SetVerified(ctx context.Context, entryCode, token string, method MethodToVerify) error
	DelVerified(ctx context.Context, entryCode, token string) error
	DecrementAttempts(ctx context.Context, entryCode, id string) (int32, error)
	CheckCooldown(ctx context.Context, entryCode, id string) time.Duration
	SetCooldown(ctx context.Context, entryCode, id string, duration time.Duration) error
	IncDailyByMethod(ctx context.Context, method *MethodToVerify) error
	IncDailyByIP(ctx context.Context, entryCode, ip string) error
}
