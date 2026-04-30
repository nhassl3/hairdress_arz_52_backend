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

func toDomainUser(user *db.User) *domain.User {
	return &domain.User{
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		FullName:    text2str(user.FullName),
		IsVerified:  user.IsVerified,
		LastLogin:   timeFromTimestampTz(user.LastLogin),
		CreatedAt:   timeFromTimestampTz(user.CreatedAt),
		UpdatedAt:   timeFromTimestampTz(user.UpdatedAt),
	}
}
