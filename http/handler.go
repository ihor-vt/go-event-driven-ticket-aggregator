package http

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"tickets/entities"
)

type Handler struct {
	eventBus    *cqrs.EventBus
	ticketsRepo TicketsRepository
	showRepo    ShowsRepository
	bookingRepo BookingRepository
}

type TicketsRepository interface {
	FindAll(ctx context.Context) ([]entities.Ticket, error)
}

type ShowsRepository interface {
	AddShow(ctx context.Context, show entities.Show) error
}

type BookingRepository interface {
	AddBooking(ctx context.Context, booking entities.Booking) error
}
