package http

import (
	"net/http"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	commandBus *cqrs.CommandBus,
	ticketsRepository TicketsRepository,
	showsRepository ShowsRepository,
	bookingsRepository BookingsRepository,
	opsBookingRepository OpsBookingRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eventBus:           eventBus,
		commandBus:         commandBus,
		ticketsRepo:        ticketsRepository,
		showsRepository:    showsRepository,
		bookingsRepository: bookingsRepository,
		opsBookingRepo:     opsBookingRepository,
	}

	e.POST("/tickets-status", handler.PostTicketsStatus)

	e.GET("/tickets", handler.GetTickets)
	e.POST("/book-tickets", handler.PostBookTickets)

	e.POST("/shows", handler.PostShows)

	e.PUT("/ticket-refund/:ticket_id", handler.PutTicketRefund)

	e.GET("/ops/bookings", handler.GetOpsBookings)
	e.GET("/ops/bookings/:id", handler.GetOpsBooking)

	return e
}
