package event

import (
	"context"
	"log/slog"

	"tickets/entities"
)

func (h Handler) CancelTicket(ctx context.Context, event entities.TicketBookingCanceled) error {
	slog.Info("Adding ticket refund to sheet")

	return h.spreadsheetsAPI.AppendRow(
		ctx,
		"tickets-to-refund",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
}
