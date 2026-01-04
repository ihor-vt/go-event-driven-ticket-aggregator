package http

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"tickets/entities"
)

type Handler struct {
	eventBus           *cqrs.EventBus
	ticketsRepo        TicketsRepository
	showsRepository    ShowsRepository
	bookingsRepository BookingsRepository
}

type TicketsRepository interface {
	FindAll(ctx context.Context) ([]entities.Ticket, error)
}

type ShowsRepository interface {
	AddShow(ctx context.Context, show entities.Show) error
}
