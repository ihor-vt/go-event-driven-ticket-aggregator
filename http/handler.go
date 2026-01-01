package http

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"tickets/entities"
)

type Handler struct {
	eventBus    *cqrs.EventBus
	ticketsRepo TicketsRepository
}

type TicketsRepository interface {
	FindAll(ctx context.Context) ([]entities.Ticket, error)
}
