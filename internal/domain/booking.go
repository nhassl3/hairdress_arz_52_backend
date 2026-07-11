package domain

import (
	"context"
	"time"
)

type BookingStatus int

const (
	UNSPECIFIED BookingStatus = iota // 0
	PENDING
	CONFIRMED
	COMPLETED
	CANCELED
	NOSHOW
)

var (
	Status = map[string]BookingStatus{
		"UNSPECIFIED": UNSPECIFIED,
		"pending":     PENDING,
		"confirmed":   CONFIRMED,
		"completed":   COMPLETED,
		"cancelled":   CANCELED,
		"no_show":     NOSHOW,
	}
)

type CreateBookingRequest struct {
	Username      string
	HairdresserID string
	ServiceID     int32
	SalonID       int32
	StartsAt      time.Time
	EndsAt        time.Time
	Description   string
	Status        BookingStatus
}

type GetBookingRequest struct {
	Username      *string
	ID            *int64
	HairdresserID *string
	ServiceID     *int32
	SalonID       *int32
}

type Booking struct {
	ID            int64
	Username      string
	HairdresserID string // UIDs
	ServiceID     int32
	SalonID       int32
	StartsAt      time.Time
	EndsAt        time.Time
	Description   string
	Status        BookingStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, params *CreateBookingRequest) (*Booking, error)
	GetBooking(ctx context.Context, params *GetBookingRequest) ([]*Booking, error)
}

// TODO: implement marshall and unmarshall JSON code for redis with encoding/json package
