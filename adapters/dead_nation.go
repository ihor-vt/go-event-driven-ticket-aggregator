package adapters

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/dead_nation"

	"tickets/entities"
)

type DeadNationClient struct {
	clients *clients.Clients
}

func NewDeadNationClient(clients *clients.Clients) *DeadNationClient {
	if clients == nil {
		panic("NewFilesApiClient: clients is nil")
	}

	return &DeadNationClient{clients: clients}
}

func (c DeadNationClient) BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error {
	resp, err := c.clients.DeadNation.PostTicketBookingWithResponse(
		ctx,
		dead_nation.PostTicketBookingRequest{
			CustomerAddress: request.CustomerEmail,
			EventId:         request.DeadNationEventID,
			NumberOfTickets: request.NumberOfTickets,
			BookingId:       request.BookingID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to book place in Dead Nation: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code from dead nation: %d", resp.StatusCode())
	}

	return nil
}
