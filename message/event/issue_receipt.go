package event

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/entities"
)

func (h Handler) IssueReceipt(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Issuing receipt")

	request := entities.IssueReceiptRequest{
		TicketID:       event.TicketID,
		Price:          event.Price,
		IdempotencyKey: event.Header.IdempotencyKey,
	}

	err := h.receiptsService.IssueReceipt(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	return nil
}
