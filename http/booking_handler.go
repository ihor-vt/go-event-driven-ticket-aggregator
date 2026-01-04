package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"tickets/entities"
)

type BookingsRepository interface {
	AddBooking(ctx context.Context, booking entities.Booking) error
}

type BookTicketRequest struct {
	CustomerEmail   string    `json:"customer_email"`
	NumberOfTickets int       `json:"number_of_tickets"`
	ShowID          uuid.UUID `json:"show_id"`
}

type BookTicketResponse struct {
	BookingID uuid.UUID `json:"booking_id"`
}

func (h Handler) PostBookTickets(c echo.Context) error {
	req := BookTicketRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.NumberOfTickets < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "number of tickets must be greater than 0")
	}

	bookingID := uuid.New()

	err := h.bookingsRepository.AddBooking(c.Request().Context(), entities.Booking{
		BookingID:       bookingID,
		CustomerEmail:   req.CustomerEmail,
		NumberOfTickets: req.NumberOfTickets,
		ShowID:          req.ShowID,
	})
	if err != nil {
		return fmt.Errorf("failed to add booking: %w", err)
	}

	return c.JSON(
		http.StatusCreated,
		BookTicketResponse{
			BookingID: bookingID,
		},
	)
}
