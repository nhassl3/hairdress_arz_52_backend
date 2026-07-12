package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nhassl3/hairdress_arz/internal/db"
	"github.com/nhassl3/hairdress_arz/internal/domain"
)

type mockDBTX struct {
	queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row
	queryFunc    func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	execFunc     func(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.queryRowFunc(ctx, sql, args...)
}

func (m *mockDBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, sql, args...)
	}
	return nil, nil
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, sql, args...)
	}
	return pgconn.CommandTag{}, nil
}

type mockRow struct {
	scanFunc func(dest ...any) error
}

func (m *mockRow) Scan(dest ...any) error {
	return m.scanFunc(dest...)
}

type mockRows struct {
	closeFunc func()
	errFunc   func() error
	nextFunc  func() bool
	scanFunc  func(dest ...any) error
	cmdTag    pgconn.CommandTag
	fields    []pgconn.FieldDescription
}

func (m *mockRows) Close() {
	if m.closeFunc != nil {
		m.closeFunc()
	}
}

func (m *mockRows) Err() error {
	if m.errFunc != nil {
		return m.errFunc()
	}
	return nil
}

func (m *mockRows) CommandTag() pgconn.CommandTag             { return m.cmdTag }
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription { return m.fields }

func (m *mockRows) Next() bool {
	if m.nextFunc != nil {
		return m.nextFunc()
	}
	return false
}

func (m *mockRows) Scan(dest ...any) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	return nil
}

func (m *mockRows) Values() ([]any, error) { return nil, nil }
func (m *mockRows) RawValues() [][]byte    { return nil }
func (m *mockRows) Conn() *pgx.Conn        { return nil }

func newMockStore(mock *mockDBTX) *db.Store {
	return &db.Store{Queries: db.New(mock)}
}

var testHairdresserUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
var testTime = time.Date(2026, 7, 1, 10, 0, 0, 0, time.UTC)

func bookingScanValues(dest ...any) error {
	*(dest[0].(*int64)) = 1
	*(dest[1].(*string)) = "testuser"
	*(dest[2].(*uuid.UUID)) = testHairdresserUUID
	*(dest[3].(*int32)) = 100
	*(dest[4].(*int32)) = 200
	*(dest[5].(*pgtype.Timestamptz)) = pgtype.Timestamptz{Time: testTime, Valid: true}
	*(dest[6].(*pgtype.Timestamptz)) = pgtype.Timestamptz{Time: testTime.Add(1 * time.Hour), Valid: true}
	*(dest[7].(*pgtype.Text)) = pgtype.Text{String: "test haircut", Valid: true}
	*(dest[8].(*string)) = "pending"
	*(dest[9].(*pgtype.Timestamptz)) = pgtype.Timestamptz{Time: testTime.Add(-1 * time.Hour), Valid: true}
	*(dest[10].(*pgtype.Timestamptz)) = pgtype.Timestamptz{Time: testTime.Add(-1 * time.Hour), Valid: true}
	return nil
}

func TestCreateBooking_Success(t *testing.T) {
	mock := &mockDBTX{
		queryRowFunc: func(_ context.Context, _ string, _ ...any) pgx.Row {
			return &mockRow{scanFunc: bookingScanValues}
		},
	}

	repo := NewBookingRepository(newMockStore(mock))

	booking, err := repo.CreateBooking(context.Background(), &domain.CreateBookingRequest{
		Username:      "testuser",
		HairdresserID: testHairdresserUUID.String(),
		ServiceID:     100,
		SalonID:       200,
		StartsAt:      testTime,
		EndsAt:        testTime.Add(1 * time.Hour),
		Description:   "test haircut",
		Status:        domain.PENDING,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if booking.ID != 1 {
		t.Errorf("ID = %d, want 1", booking.ID)
	}
	if booking.Username != "testuser" {
		t.Errorf("Username = %q, want %q", booking.Username, "testuser")
	}
	if booking.HairdresserID != testHairdresserUUID.String() {
		t.Errorf("HairdresserID = %q, want %q", booking.HairdresserID, testHairdresserUUID.String())
	}
	if booking.ServiceID != 100 {
		t.Errorf("ServiceID = %d, want 100", booking.ServiceID)
	}
	if booking.SalonID != 200 {
		t.Errorf("SalonID = %d, want 200", booking.SalonID)
	}
	if !booking.StartsAt.Equal(testTime) {
		t.Errorf("StartsAt = %v, want %v", booking.StartsAt, testTime)
	}
	if !booking.EndsAt.Equal(testTime.Add(1 * time.Hour)) {
		t.Errorf("EndsAt = %v, want %v", booking.EndsAt, testTime.Add(1*time.Hour))
	}
	if booking.Description != "test haircut" {
		t.Errorf("Description = %q, want %q", booking.Description, "test haircut")
	}
	if booking.Status != domain.PENDING {
		t.Errorf("Status = %v, want %v", booking.Status, domain.PENDING)
	}
}

func TestCreateBooking_PgError23514_ReturnsErrDataNoProvide(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23514", Message: "violates check constraint"}

	mock := &mockDBTX{
		queryRowFunc: func(_ context.Context, _ string, _ ...any) pgx.Row {
			return &mockRow{scanFunc: func(_ ...any) error { return pgErr }}
		},
	}

	repo := NewBookingRepository(newMockStore(mock))

	_, err := repo.CreateBooking(context.Background(), &domain.CreateBookingRequest{
		Username:      "testuser",
		HairdresserID: testHairdresserUUID.String(),
		ServiceID:     100,
		SalonID:       200,
		StartsAt:      testTime,
		EndsAt:        testTime.Add(1 * time.Hour),
		Description:   "test",
		Status:        domain.PENDING,
	})
	if !errors.Is(err, domain.ErrDataNoProvide) {
		t.Errorf("expected ErrDataNoProvide, got %v", err)
	}
}

func TestCreateBooking_OtherPgError_WrapsWithCode(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505", Message: "unique violation"}

	mock := &mockDBTX{
		queryRowFunc: func(_ context.Context, _ string, _ ...any) pgx.Row {
			return &mockRow{scanFunc: func(_ ...any) error { return pgErr }}
		},
	}

	repo := NewBookingRepository(newMockStore(mock))

	_, err := repo.CreateBooking(context.Background(), &domain.CreateBookingRequest{
		Username:      "testuser",
		HairdresserID: testHairdresserUUID.String(),
		ServiceID:     100,
		SalonID:       200,
		StartsAt:      testTime,
		EndsAt:        testTime.Add(1 * time.Hour),
		Description:   "test",
		Status:        domain.PENDING,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.As(err, &pgErr) {
		t.Error("expected error to be a PgError")
	}
}

func TestCreateBooking_GenericError_Wraps(t *testing.T) {
	genericErr := errors.New("connection refused")

	mock := &mockDBTX{
		queryRowFunc: func(_ context.Context, _ string, _ ...any) pgx.Row {
			return &mockRow{scanFunc: func(_ ...any) error { return genericErr }}
		},
	}

	repo := NewBookingRepository(newMockStore(mock))

	_, err := repo.CreateBooking(context.Background(), &domain.CreateBookingRequest{
		Username:      "testuser",
		HairdresserID: testHairdresserUUID.String(),
		ServiceID:     100,
		SalonID:       200,
		StartsAt:      testTime,
		EndsAt:        testTime.Add(1 * time.Hour),
		Description:   "test",
		Status:        domain.PENDING,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, genericErr) {
		t.Error("expected error to wrap the original error")
	}
}

func TestGetBooking_Success(t *testing.T) {
	var called bool

	mock := &mockDBTX{
		queryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &mockRows{
				nextFunc: func() bool {
					if !called {
						called = true
						return true
					}
					return false
				},
				scanFunc:  bookingScanValues,
				closeFunc: func() {},
				errFunc:   func() error { return nil },
			}, nil
		},
	}

	repo := NewBookingRepository(newMockStore(mock))

	username := "testuser"
	bookings, err := repo.GetBooking(context.Background(), &domain.GetBookingRequest{
		Username: &username,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bookings) == 0 {
		t.Fatal("expected at least one booking")
	}

	var booking *domain.Booking
	for _, b := range bookings {
		if b != nil {
			booking = b
			break
		}
	}
	if booking == nil {
		t.Fatal("expected a non-nil booking in results")
	}
	if booking.ID != 1 {
		t.Errorf("ID = %d, want 1", booking.ID)
	}
	if booking.Username != "testuser" {
		t.Errorf("Username = %q, want %q", booking.Username, "testuser")
	}
	if booking.Status != domain.PENDING {
		t.Errorf("Status = %v, want %v", booking.Status, domain.PENDING)
	}
}

func TestGetBooking_Empty_ReturnsErrNoBookings(t *testing.T) {
	mock := &mockDBTX{
		queryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return &mockRows{
				nextFunc:  func() bool { return false },
				closeFunc: func() {},
				errFunc:   func() error { return nil },
			}, nil
		},
	}

	repo := NewBookingRepository(newMockStore(mock))

	username := "nonexistent"
	_, err := repo.GetBooking(context.Background(), &domain.GetBookingRequest{
		Username: &username,
	})
	if !errors.Is(err, domain.ErrNoBookings) {
		t.Errorf("expected ErrNoBookings, got %v", err)
	}
}

func TestGetBooking_QueryError_Wraps(t *testing.T) {
	queryErr := errors.New("query failed")

	mock := &mockDBTX{
		queryFunc: func(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
			return nil, queryErr
		},
	}

	repo := NewBookingRepository(newMockStore(mock))

	username := "testuser"
	_, err := repo.GetBooking(context.Background(), &domain.GetBookingRequest{
		Username: &username,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, queryErr) {
		t.Error("expected error to wrap the original query error")
	}
}
