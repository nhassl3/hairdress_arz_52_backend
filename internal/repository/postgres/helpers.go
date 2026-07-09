package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhassl3/hairdress_arz/internal/db"
	"github.com/nhassl3/hairdress_arz/internal/domain"
)

func str2Text(str *string) pgtype.Text {
	if str != nil && *str != "" {
		return pgtype.Text{
			String: *str,
			Valid:  true,
		}
	}
	return pgtype.Text{}
}

func timeFromTimestampTz(tm pgtype.Timestamptz) time.Time {
	if tm.Valid {
		return tm.Time
	}
	return time.Time{}
}

func text2str(text pgtype.Text) string {
	return text.String
}

func pgTimeTZ(ts pgtype.Timestamptz, _ *time.Location) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}

func toDomainUser(user *db.User) *domain.User {
	return &domain.User{
		Username:    user.Username,
		UID:         user.Uid.String(),
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		FullName:    text2str(user.FullName),
		IsVerified:  user.IsVerified,
		LastLogin:   timeFromTimestampTz(user.LastLogin),
		CreatedAt:   timeFromTimestampTz(user.CreatedAt),
		UpdatedAt:   timeFromTimestampTz(user.UpdatedAt),
	}
}

func toDomainSession(s db.Session) *domain.Session {
	return &domain.Session{
		ID:           s.ID.String(),
		Username:     s.Username,
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		ClientIP:     s.ClientIp,
		ExpiresAt:    pgTimeTZ(s.ExpiresAt, time.UTC),
		IsBlocked:    s.IsBlocked,
		CreatedAt:    pgTimeTZ(s.CreatedAt, time.UTC),
	}
}
