package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhassl3/hairdress_arz/internal/db"
	"github.com/nhassl3/hairdress_arz/internal/domain"
)

type BookingRepository struct {
	store *db.Store
}

func NewBookingRepository(store *db.Store) *BookingRepository {
	return &BookingRepository{
		store: store,
	}
}

// CreateBooking creates booking record in DB. Status by default pending
func (repo *BookingRepository) CreateBooking(ctx context.Context, params *domain.CreateBookingRequest) (*domain.Booking, error) {
	hairdresserUUID, err := string2UUID(params.HairdresserID)
	if err != nil {
		return nil, fmt.Errorf("booking_repo.CreateBooking: invalid hairdresser_id: %w", domain.ErrDataNoProvide)
	}

	booking, err := repo.store.CreateBooking(ctx, db.CreateBookingParams{
		Username:      params.Username,
		HairdresserID: hairdresserUUID,
		ServiceID:     params.ServiceID,
		SalonID:       params.SalonID,
		StartsAt:      time2PgTime(params.StartsAt),
		EndsAt:        time2PgTime(params.EndsAt),
		Description:   str2Text(&params.Description),
		Status:        domain.Status[domain.PENDING], // an explicit indication that the status is "pending"
	})
	if err != nil {
		return nil, handleErrors(err, "CreateBooking: failed to create booking")
	}
	return toDomainBooking(booking), nil
}

// GetBookings gets record from DB about booking by several parameters
// - username
// - ID
// - Hairdresser UID
// - Service ID
// - Salon ID
// All parameters except ID also accept the start time to search for relevant records.
func (repo *BookingRepository) GetBookings(ctx context.Context, params *domain.GetBookingRequest) ([]*domain.Booking, error) {
	hairdresserUUID, err := string2PgUUID(params.HairdresserID)
	if err != nil {
		return nil, fmt.Errorf("booking_repo.GetBooking: invalid hairdresser_id: %w", domain.ErrDataNoProvide)
	}

	bookings, err := repo.store.GetBooking(ctx, db.GetBookingParams{
		Username:      str2Text(params.Username),
		ID:            int2Int8(params.ID),
		HairdresserID: hairdresserUUID,
		ServiceID:     int2Int4(params.ServiceID),
		SalonID:       int2Int4(params.SalonID),
	})
	if err != nil {
		return nil, handleErrors(err, "GetBookings: failed to get bookings")
	}
	domainBookings := toDomainBookings(bookings)
	if domainBookings == nil {
		return nil, domain.ErrNoBookings
	}
	return domainBookings, nil
}

// UpdateBookingStatus updates status in record of user booking
// status can be unspecified (system zero value), pending (default value when creating), confirmed (checked and approved by owner)
// completed (mark like a successfully complete hairdress work) and canceled
func (repo *BookingRepository) UpdateBookingStatus(ctx context.Context, params *domain.UpdateBookingStatusRequest) (*domain.Booking, error) {
	var (
		targetTime     time.Time
		hairdresserUID pgtype.UUID
		err            error
	)

	if params.FindByUsername != nil && (params.FindByUsername.StartsAt != time.Time{}) {
		targetTime = params.FindByUsername.StartsAt
	} else if params.FindByService != nil && (params.FindByService.StartsAt != time.Time{}) {
		targetTime = params.FindByService.StartsAt
	} else if params.FindBySalon != nil && (params.FindBySalon.StartsAt != time.Time{}) {
		targetTime = params.FindBySalon.StartsAt
	} else if params.FindByHairdresser != nil && (params.FindByHairdresser.StartsAt != time.Time{}) {
		targetTime = params.FindByHairdresser.StartsAt
	} else {
		targetTime = time.Now().Truncate(15 * time.Minute) // -15 minutes before now
	}

	if params.FindByHairdresser != nil {
		hairdresserUID, err = string2PgUUID(&params.FindByHairdresser.HairdresserUID)
		if err != nil {
			return nil, fmt.Errorf("booking_repo.UpdateBookingStatus: invalid hairdresser_id: %w", domain.ErrDataNoProvide)
		}
	}

	booking, err := repo.store.UpdateBookingStatus(ctx, db.UpdateBookingStatusParams{
		BookingStatus:  domain.Status[params.Status],
		ID:             int2Int8(params.ID),
		Username:       str2Text(&params.FindByUsername.Username),
		ServiceID:      int2Int4(&params.FindByService.ServiceID),
		SalonID:        int2Int4(&params.FindBySalon.SalonID),
		HairdresserUid: hairdresserUID,
		TargetTime:     time2PgTime(targetTime),
	})
	if err != nil {
		return nil, handleErrors(err, "UpdateBookingStatus: failed to update booking")
	}

	domainBooking := toDomainBooking(booking)
	if domainBooking == nil {
		return nil, domain.ErrNoBookings
	}

	return domainBooking, nil
}

func handleErrors(err error, defaultMessage string) error {
	if err == nil {
		return nil
	}
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23514":
			return domain.ErrDataNoProvide
		case "23503":
			return domain.ErrFailedToFindInTheSystem
		default:
			return fmt.Errorf("booking_repo: %s: pg err. CODE: %s; %w", defaultMessage, strings.ToUpper(pgErr.Code), err)
		}
	}
	return fmt.Errorf("booking_repo: %s: %w", defaultMessage, err)
}
