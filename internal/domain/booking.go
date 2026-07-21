package domain

import (
	"context"
	"time"

	reverseEnums "github.com/nhassl3/hairdress_arz/pkg/reverse-enums"
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
	Status = map[BookingStatus]string{
		UNSPECIFIED: "UNSPECIFIED",
		PENDING:     "pending",
		CONFIRMED:   "confirmed",
		COMPLETED:   "completed",
		CANCELED:    "cancelled",
		NOSHOW:      "no_show",
	}
	StatusReversed = reverseEnums.ReverseMap(Status)
)

type CreateBookingRequest struct {
	Username      string
	HairdresserID string
	ServiceID     int32
	SalonID       int32
	StartsAt      time.Time
	EndsAt        time.Time
	Description   string
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
	GetBookings(ctx context.Context, params *GetBookingRequest) ([]*Booking, error)
	UpdateBookingStatus(ctx context.Context, params *UpdateBookingStatusRequest) (*Booking, error)
}

type FindByUsername struct {
	Username string
	StartsAt time.Time
}

type FindByHairdresser struct {
	HairdresserUID string
	StartsAt       time.Time
}

type FindByService struct {
	ServiceID int32
	StartsAt  time.Time
}

type FindBySalon struct {
	SalonID  int32
	StartsAt time.Time
}

type UpdateBookingStatusRequest struct {
	ID                *int64
	FindByUsername    *FindByUsername
	FindByHairdresser *FindByHairdresser
	FindByService     *FindByService
	FindBySalon       *FindBySalon
	Status            BookingStatus
}

func NewUpdateBookingStatusRequest(status BookingStatus) *UpdateBookingStatusRequest {
	return &UpdateBookingStatusRequest{
		Status: status,
	}
}

func (upd *UpdateBookingStatusRequest) WithID(ID *int64) *UpdateBookingStatusRequest {
	return &UpdateBookingStatusRequest{
		ID:     ID,
		Status: upd.Status,
	}
}

func (upd *UpdateBookingStatusRequest) WithUsername(username string, startsAt time.Time) *UpdateBookingStatusRequest {
	return &UpdateBookingStatusRequest{
		FindByUsername: &FindByUsername{
			Username: username,
			StartsAt: startsAt,
		},
		Status: upd.Status,
	}
}

func (upd *UpdateBookingStatusRequest) WithService(serviceID int32, startsAt time.Time) *UpdateBookingStatusRequest {
	return &UpdateBookingStatusRequest{
		FindByService: &FindByService{
			ServiceID: serviceID,
			StartsAt:  startsAt,
		},
		Status: upd.Status,
	}
}

func (upd *UpdateBookingStatusRequest) WithSalon(salonID int32, startsAt time.Time) *UpdateBookingStatusRequest {
	return &UpdateBookingStatusRequest{
		FindBySalon: &FindBySalon{
			SalonID:  salonID,
			StartsAt: startsAt,
		},
		Status: upd.Status,
	}
}

func (upd *UpdateBookingStatusRequest) WithHairdresser(hairdresserUID string, startsAt time.Time) *UpdateBookingStatusRequest {
	return &UpdateBookingStatusRequest{
		FindByHairdresser: &FindByHairdresser{
			HairdresserUID: hairdresserUID,
			StartsAt:       startsAt,
		},
		Status: upd.Status,
	}
}

// TODO: implement marshall and unmarshall JSON code for redis with encoding/json package
