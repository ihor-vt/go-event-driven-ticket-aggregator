package event

import (
	"context"
	"log/slog"

	"tickets/entities"
)

func (h Handler) AppendToTracker(ctx context.Context, event entities.TicketBookingConfirmed) error {
	slog.Info("Appending ticket to the tracker")

	return h.spreadsheetsAPI.AppendRow(
		ctx,
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
}
