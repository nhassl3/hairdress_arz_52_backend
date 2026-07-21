package grpc

import (
	"context"
	"fmt"

	"github.com/nhassl3/hairdress_arz/internal/domain"
	"github.com/nhassl3/hairdress_arz/internal/service"
	reverseEnums "github.com/nhassl3/hairdress_arz/pkg/reverse-enums"
	bookingv1 "github.com/nhassl3/hairdress_arz_52_contracts/pkg/pb/booking/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	StatusBooking = map[bookingv1.StatusBooking]domain.BookingStatus{
		bookingv1.StatusBooking_BOOKING_STATUS_UNSPECIFIED: domain.UNSPECIFIED,
		bookingv1.StatusBooking_BOOKING_STATUS_PENDING:     domain.PENDING,
		bookingv1.StatusBooking_BOOKING_STATUS_COMPLETED:   domain.COMPLETED,
		bookingv1.StatusBooking_BOOKING_STATUS_CONFIRMED:   domain.CONFIRMED,
		bookingv1.StatusBooking_BOOKING_STATUS_CANCELLED:   domain.CANCELED,
		bookingv1.StatusBooking_BOOKING_STATUS_NO_SHOW:     domain.NOSHOW,
	}
	StatusBookingReversed = reverseEnums.ReverseMap(StatusBooking)
)

type BookingHandler struct {
	bookingv1.UnimplementedBookingServiceServer
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{
		svc: svc,
	}
}

func (h *BookingHandler) CreateBooking(ctx context.Context, req *bookingv1.CreateBookingRequest) (*bookingv1.CreateBookingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	booking, err := h.svc.CreateBooking(ctx, &domain.CreateBookingRequest{
		Username:      req.Username,
		HairdresserID: req.HairdresserId,
		ServiceID:     req.ServiceId,
		SalonID:       req.SalonId,
		Description:   req.Description,
		StartsAt:      req.StartsAt.AsTime(),
		EndsAt:        req.EndsAt.AsTime(),
	})
	if err != nil {
		return nil, domainErr(err)
	}
	return &bookingv1.CreateBookingResponse{
		Booking: mapBooking(booking),
	}, nil
}

func (h *BookingHandler) GetBookings(ctx context.Context, req *bookingv1.GetBookingRequest) (*bookingv1.GetBookingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	var (
		bookings []*domain.Booking
		err      error
	)
	switch v := req.GetMethod().(type) {
	case *bookingv1.GetBookingRequest_Username:
		bookings, err = h.svc.GetBookings(ctx, &domain.GetBookingRequest{
			Username: &v.Username,
		})
	case *bookingv1.GetBookingRequest_ServiceId:
		bookings, err = h.svc.GetBookings(ctx, &domain.GetBookingRequest{
			ServiceID: &v.ServiceId,
		})
	case *bookingv1.GetBookingRequest_SalonId:
		bookings, err = h.svc.GetBookings(ctx, &domain.GetBookingRequest{
			SalonID: &v.SalonId,
		})
	case *bookingv1.GetBookingRequest_HairdresserId:
		bookings, err = h.svc.GetBookings(ctx, &domain.GetBookingRequest{
			HairdresserID: &v.HairdresserId,
		})
	case *bookingv1.GetBookingRequest_Id:
		bookings, err = h.svc.GetBookings(ctx, &domain.GetBookingRequest{
			ID: &v.Id,
		})
	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("%T(%v)", v, v))
	}
	if err != nil {
		return nil, domainErr(err)
	}
	mBookings := mapBookings(bookings)
	if mBookings == nil {
		return nil, domainErr(domain.ErrNoBookings)
	}
	return &bookingv1.GetBookingResponse{
		Bookings: mBookings,
	}, nil
}

func (h *BookingHandler) UpdateBookingStatus(ctx context.Context, req *bookingv1.UpdateBookingStatusRequest) (*bookingv1.UpdateBookingStatusResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var (
		booking *domain.Booking
		err     error
		params  = domain.NewUpdateBookingStatusRequest(StatusBooking[req.GetNewStatus()])
	)
	switch v := req.GetMethod().(type) {
	case *bookingv1.UpdateBookingStatusRequest_Id:
		booking, err = h.svc.UpdateBookingStatus(ctx, params.WithID(&v.Id))
	case *bookingv1.UpdateBookingStatusRequest_FindByHairdresser:
		booking, err = h.svc.UpdateBookingStatus(
			ctx, params.WithHairdresser(v.FindByHairdresser.GetHairdresserId(), v.FindByHairdresser.GetStartsAt().AsTime()))
	case *bookingv1.UpdateBookingStatusRequest_FindByService:
		booking, err = h.svc.UpdateBookingStatus(ctx,
			params.WithService(v.FindByService.GetServiceId(), v.FindByService.GetStartsAt().AsTime()))
	case *bookingv1.UpdateBookingStatusRequest_FindBySalon:
		booking, err = h.svc.UpdateBookingStatus(ctx,
			params.WithSalon(v.FindBySalon.GetSalonId(), v.FindBySalon.GetStartsAt().AsTime()))
	case *bookingv1.UpdateBookingStatusRequest_FindByUsername:
		booking, err = h.svc.UpdateBookingStatus(ctx,
			params.WithUsername(v.FindByUsername.GetUsername(), v.FindByUsername.GetStartsAt().AsTime()))
	default:
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("%T(%v)", v, v))
	}
	if err != nil {
		return nil, domainErr(err)
	}

	mBooking := mapBooking(booking)
	if mBooking == nil {
		return nil, domainErr(domain.ErrNoBookings)
	}

	return &bookingv1.UpdateBookingStatusResponse{
		Booking: mBooking,
	}, nil
}

func mapBooking(booking *domain.Booking) *bookingv1.Bookings {
	if booking == nil {
		return nil
	}
	return &bookingv1.Bookings{
		Id:            booking.ID,
		Username:      booking.Username,
		HairdresserId: booking.HairdresserID,
		ServiceId:     booking.ServiceID,
		SalonId:       booking.SalonID,
		Description:   booking.Description,
		StartsAt:      safeTimestamp(booking.StartsAt),
		EndsAt:        safeTimestamp(booking.EndsAt),
		Status:        StatusBookingReversed[booking.Status],
		CreatedAt:     safeTimestamp(booking.CreatedAt),
		UpdatedAt:     safeTimestamp(booking.UpdatedAt),
	}
}

func mapBookings(bookings []*domain.Booking) []*bookingv1.Bookings {
	zap.L().Info("INFO", zap.Any("bookings", bookings))
	if bookings == nil || len(bookings) == 0 {
		return nil
	}
	protoBookings := make([]*bookingv1.Bookings, 0, len(bookings))
	for _, booking := range bookings {
		protoBooking := mapBooking(booking)
		if protoBooking != nil {
			protoBookings = append(protoBookings, protoBooking)
		}
	}
	return protoBookings
}
