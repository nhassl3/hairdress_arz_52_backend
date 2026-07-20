package postgres

import (
	"reflect"
	"time"

	"github.com/google/uuid"
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

func int2Int8(i *int64) pgtype.Int8 {
	if i == nil || *i == 0 {
		return pgtype.Int8{}
	}
	return pgtype.Int8{
		Int64: *i,
		Valid: true,
	}
}

func int2Int4(i *int32) pgtype.Int4 {
	if i == nil || *i == 0 {
		return pgtype.Int4{}
	}
	return pgtype.Int4{
		Int32: *i,
		Valid: true,
	}
}

func string2UUID(str string) (uuid.UUID, error) {
	if str == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(str)
}

func string2PgUUID(str *string) (pgtype.UUID, error) {
	if str == nil || len(*str) == 0 {
		return pgtype.UUID{}, nil
	}
	id, err := uuid.Parse(*str)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: id, Valid: true}, nil
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

func time2PgTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
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

func toDomainBooking(booking db.Booking) *domain.Booking {
	if reflect.ValueOf(booking).IsZero() {
		return nil
	}
	return &domain.Booking{
		ID:            booking.ID,
		Username:      booking.Username,
		HairdresserID: booking.HairdresserID.String(),
		ServiceID:     booking.ServiceID,
		SalonID:       booking.SalonID,
		StartsAt:      pgTimeTZ(booking.StartsAt, time.UTC),
		EndsAt:        pgTimeTZ(booking.EndsAt, time.UTC),
		Description:   text2str(booking.Description),
		Status:        domain.StatusReversed[booking.Status],
		CreatedAt:     pgTimeTZ(booking.CreatedAt, time.UTC),
		UpdatedAt:     pgTimeTZ(booking.UpdatedAt, time.UTC),
	}
}

func toDomainBookings(bookings []db.Booking) []*domain.Booking {
	if bookings == nil || len(bookings) == 0 {
		return nil
	}
	domainBookings := make([]*domain.Booking, len(bookings))
	for _, booking := range bookings {
		mapBooking := toDomainBooking(booking)
		if mapBooking != nil {
			domainBookings = append(domainBookings, mapBooking)
		}
	}
	return domainBookings
}
