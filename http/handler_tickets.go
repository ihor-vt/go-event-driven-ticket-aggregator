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

	idempotencyKey := c.Request().Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Idempotency-Key header is required")
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			event := entities.TicketBookingConfirmed{
				Header: entities.NewMessageHeaderWithIdempotencyKey(idempotencyKey + ticket.TicketID),

				TicketID:      ticket.TicketID,
				Price:         ticket.Price,
				CustomerEmail: ticket.CustomerEmail,
			}

			if err := h.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("failed to publish TicketBookingConfirmed event: %w", err)
			}
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewMessageHeaderWithIdempotencyKey(idempotencyKey + ticket.TicketID),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			if err := h.eventBus.Publish(c.Request().Context(), event); err != nil {
				return fmt.Errorf("failed to publish TicketBookingCanceled event: %w", err)
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) GetTickets(c echo.Context) error {
	tickets, err := h.ticketsRepo.FindAll(c.Request().Context())
	if err != nil {
		return fmt.Errorf("failed to find tickets: %w", err)
	}

	return c.JSON(http.StatusOK, tickets)
}
