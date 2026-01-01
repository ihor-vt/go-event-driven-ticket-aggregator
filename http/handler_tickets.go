package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"tickets/entities"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatusRequest `json:"tickets"`
}

type TicketStatusRequest struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
}

func (h Handler) PostTicketsStatus(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			event := entities.TicketBookingConfirmed{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			if err := h.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("publishing TicketBookingConfirmed event: %w", err)
			}

			// msg := message.NewMessage(watermill.NewUUID(), payload)
			// msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))
			// msg.Metadata.Set("type", "TicketBookingConfirmed")

			// err = h.publisher.Publish("TicketBookingConfirmed", msg)
			// if err != nil {
			// 	return err
			// }
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			if err := h.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("publishing TicketBookingCanceled event: %w", err)
			}

			// msg := message.NewMessage(watermill.NewUUID(), payload)
			// msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))
			// msg.Metadata.Set("type", "TicketBookingCanceled")

			// err = h.publisher.Publish("TicketBookingCanceled", msg)
			// if err != nil {
			// 	return err
			// }
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}
