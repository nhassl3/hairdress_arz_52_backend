package service

import (
	"context"
	"fmt"

	"github.com/nhassl3/hairdress_arz/internal/domain"
)

type BookingService struct {
	bookingRepo domain.BookingRepository
}

func NewBookingService(bookingRepo domain.BookingRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo}
}

func (service *BookingService) CreateBooking(ctx context.Context, params *domain.CreateBookingRequest) (*domain.Booking, error) {
	booking, err := service.bookingRepo.CreateBooking(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("booking_service.CreateBooking: failed to create booking: %w", err)
	}
	return booking, nil
}

func (service *BookingService) GetBookings(ctx context.Context, params *domain.GetBookingRequest) ([]*domain.Booking, error) {
	bookings, err := service.bookingRepo.GetBookings(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("booking_service.GetBooking: failed to get bookings: %w", err)
	}
	return bookings, nil
}

func (service *BookingService) UpdateBookingStatus(ctx context.Context, params *domain.UpdateBookingStatusRequest) (*domain.Booking, error) {
	booking, err := service.bookingRepo.UpdateBookingStatus(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("booking_service.UpdateBookingStatus: failed to patch booking status: %w", err)
	}
	return booking, nil
}
