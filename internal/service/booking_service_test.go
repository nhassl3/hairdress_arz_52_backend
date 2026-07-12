package service

import (
	"context"
	"errors"
	"testing"

	"github.com/nhassl3/hairdress_arz/internal/domain"
)

type mockBookingRepo struct {
	createBookingFunc func(ctx context.Context, params *domain.CreateBookingRequest) (*domain.Booking, error)
	getBookingFunc    func(ctx context.Context, params *domain.GetBookingRequest) ([]*domain.Booking, error)
}

func (m *mockBookingRepo) CreateBooking(ctx context.Context, params *domain.CreateBookingRequest) (*domain.Booking, error) {
	return m.createBookingFunc(ctx, params)
}

func (m *mockBookingRepo) GetBooking(ctx context.Context, params *domain.GetBookingRequest) ([]*domain.Booking, error) {
	return m.getBookingFunc(ctx, params)
}

func TestCreateBooking_Success(t *testing.T) {
	expected := &domain.Booking{
		ID:        1,
		Username:  "testuser",
		ServiceID: 100,
		SalonID:   200,
		Status:    domain.PENDING,
	}

	repo := &mockBookingRepo{
		createBookingFunc: func(_ context.Context, params *domain.CreateBookingRequest) (*domain.Booking, error) {
			if params == nil {
				t.Error("params should not be nil")
			}
			return expected, nil
		},
	}

	svc := NewBookingService(repo)

	booking, err := svc.CreateBooking(context.Background(), &domain.CreateBookingRequest{
		Username:  "testuser",
		ServiceID: 100,
		SalonID:   200,
		Status:    domain.PENDING,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if booking != expected {
		t.Errorf("returned booking != expected")
	}
}

func TestCreateBooking_RepoError_Wraps(t *testing.T) {
	repoErr := errors.New("db error")

	repo := &mockBookingRepo{
		createBookingFunc: func(_ context.Context, _ *domain.CreateBookingRequest) (*domain.Booking, error) {
			return nil, repoErr
		},
	}

	svc := NewBookingService(repo)

	_, err := svc.CreateBooking(context.Background(), &domain.CreateBookingRequest{
		Username:  "testuser",
		ServiceID: 100,
		SalonID:   200,
		Status:    domain.PENDING,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Error("expected error to wrap the repository error")
	}
}

func TestGetBookings_Success(t *testing.T) {
	expected := []*domain.Booking{
		{ID: 1, Username: "testuser", Status: domain.PENDING},
		{ID: 2, Username: "testuser2", Status: domain.CONFIRMED},
	}

	repo := &mockBookingRepo{
		getBookingFunc: func(_ context.Context, params *domain.GetBookingRequest) ([]*domain.Booking, error) {
			if params == nil {
				t.Error("params should not be nil")
			}
			return expected, nil
		},
	}

	svc := NewBookingService(repo)

	username := "testuser"
	bookings, err := svc.GetBookings(context.Background(), &domain.GetBookingRequest{
		Username: &username,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bookings) != 2 {
		t.Fatalf("expected 2 bookings, got %d", len(bookings))
	}
	for i := range bookings {
		if bookings[i] != expected[i] {
			t.Errorf("bookings[%d] mismatch", i)
		}
	}
}

func TestGetBookings_RepoError_Wraps(t *testing.T) {
	repoErr := errors.New("db error")

	repo := &mockBookingRepo{
		getBookingFunc: func(_ context.Context, _ *domain.GetBookingRequest) ([]*domain.Booking, error) {
			return nil, repoErr
		},
	}

	svc := NewBookingService(repo)

	username := "testuser"
	_, err := svc.GetBookings(context.Background(), &domain.GetBookingRequest{
		Username: &username,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repoErr) {
		t.Error("expected error to wrap the repository error")
	}
}
