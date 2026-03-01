package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetOpsBookings(c echo.Context) error {
	bookings, err := h.opsBookingRepo.AllBookings(c.Request().Context())
	if err != nil {
		return fmt.Errorf("failed to get bookings: %w", err)
	}

	return c.JSON(http.StatusOK, bookings)
}

func (h Handler) GetOpsBooking(c echo.Context) error {
	bookingID := c.Param("id")

	booking, err := h.opsBookingRepo.BookingReadModel(c.Request().Context(), bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	return c.JSON(http.StatusOK, booking)
}
