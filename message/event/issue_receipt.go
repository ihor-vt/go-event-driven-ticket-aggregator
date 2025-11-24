package event

import (
	"context"
	"fmt"
	"log/slog"

	"tickets/entities"
)

func (h Handler) IssueReceipt(ctx context.Context, event entities.TicketBookingConfirmed) error {
	slog.Info("Issuing receipt")

	request := entities.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price:    event.Price,
	}

	err := h.receiptsService.IssueReceipt(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	return nil
}
