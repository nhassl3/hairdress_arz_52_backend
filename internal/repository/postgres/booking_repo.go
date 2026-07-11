package postgres

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
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

func (repo *BookingRepository) CreateBooking(ctx context.Context, params *domain.CreateBookingRequest) (*domain.Booking, error) {
	booking, err := repo.store.CreateBooking(ctx, db.CreateBookingParams{
		Username:      params.Username,
		HairdresserID: string2UUID(params.HairdresserID),
		ServiceID:     params.ServiceID,
		SalonID:       params.SalonID,
		StartsAt:      time2PgTime(params.StartsAt),
		EndsAt:        time2PgTime(params.EndsAt),
		Description:   str2Text(&params.Description),
		Status:        slices.Sorted(maps.Keys(domain.Status))[params.Status],
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case "23514":
				return nil, domain.ErrDataNoProvide
			default:
				return nil, fmt.Errorf("booking_repo.CreateBooking: pg err. CODE: %s; %w", strings.ToUpper(pgErr.Code), err)
			}
		}
		return nil, fmt.Errorf("booking_repo.CreateBooking: failed to create booking: %w", err)
	}
	return toDomainBooking(booking), nil
}

func (repo *BookingRepository) GetBooking(ctx context.Context, params *domain.GetBookingRequest) ([]*domain.Booking, error) {
	bookings, err := repo.store.GetBooking(ctx, db.GetBookingParams{
		Username:      str2Text(params.Username),
		ID:            int2Int8(params.ID),
		HairdresserID: string2PgUUID(params.HairdresserID),
		ServiceID:     int2Int4(params.ServiceID),
		SalonID:       int2Int4(params.SalonID),
	})
	if err != nil {
		return nil, fmt.Errorf("booking_repo.GetBooking: failed to get bookings: %w", err)
	}
	domainBookings := toDomainBookings(bookings)
	if domainBookings == nil {
		return nil, domain.ErrNoBookings
	}
	return domainBookings, nil
}
