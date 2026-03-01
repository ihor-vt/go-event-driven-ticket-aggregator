package http

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"tickets/entities"
)

type Handler struct {
	eventBus           *cqrs.EventBus
	commandBus         *cqrs.CommandBus
	ticketsRepo        TicketsRepository
	showsRepository    ShowsRepository
	bookingsRepository BookingsRepository
	opsBookingRepo     OpsBookingRepository
}

type TicketsRepository interface {
	FindAll(ctx context.Context) ([]entities.Ticket, error)
}

type ShowsRepository interface {
	AddShow(ctx context.Context, show entities.Show) error
}

type OpsBookingRepository interface {
	AllBookings(ctx context.Context) ([]entities.OpsBooking, error)
	BookingReadModel(ctx context.Context, bookingID string) (entities.OpsBooking, error)
}
