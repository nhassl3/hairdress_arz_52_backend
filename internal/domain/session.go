package domain

import (
	"encoding/json"
	"time"
)

type Session struct {
	ID           string    `db:"id"`
	Username     string    `db:"username"`
	RefreshToken string    `db:"refresh_token"`
	UserAgent    string    `db:"user_agent"`
	ClientIP     string    `db:"client_ip"`
	IsBlocked    bool      `db:"is_blocked"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
}

// MarshalBinary this method needed for correct work of Redis
// because Redis only work with JSON, but not with a structures
// marshalling Session structure to the []byte code (JSON)
func (u *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// UnmarshalBinary this method needed for correct work of Redis
// because Redis only work with JSON, but not with a structures
// unmarshalling source data and convert to Session structure
func (u *Session) UnmarshalBinary(data []byte) error {
	if u == nil {
		return ErrRedisNotFound
	}
	return json.Unmarshal(data, u)
}

type CreateSessionParams struct {
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIp     string
	IsBlocked    bool
	ExpiresAt    time.Time
}
