package event

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/entities"
)

func (h Handler) TicketRefundToSheet(ctx context.Context, event *entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Info("Adding ticket refund to sheet")

	return h.spreadsheetsAPI.AppendRow(
		ctx,
		"tickets-to-refund",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
}
