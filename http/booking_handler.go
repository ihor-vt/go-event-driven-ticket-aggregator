package http

import (
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BookTicketsRequest struct {
	ShowID          uuid.UUID `json:"show_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email"`
}

func (h Handler) PostBookTickets(c echo.Context) error {
	var req BookTicketsRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	bookingId := uuid.New()

	if err := h.bookingRepo.AddBooking(c.Request().Context(), entities.Booking{
		BookingID:       bookingId.String(),
		ShowID:          req.ShowID.String(),
		NumberOfTickets: req.NumberOfTickets,
		CustomerEmail:   req.CustomerEmail,
	}); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"booking_id": bookingId.String(),
	})
}
