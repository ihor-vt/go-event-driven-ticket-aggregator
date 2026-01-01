package event

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/entities"
)

func (h Handler) RemoveCanceledTicket(ctx context.Context, event *entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Info("Storing ticket")

	return h.ticketsRepository.Remove(ctx, event.TicketID)
}
