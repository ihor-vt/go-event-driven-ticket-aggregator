package http

import (
	"net/http"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	ticketsRepository TicketsRepository,
	showRepo ShowsRepository,
	bookingRepo BookingRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eventBus:    eventBus,
		ticketsRepo: ticketsRepository,
		showRepo:    showRepo,
		bookingRepo: bookingRepo,
	}

	e.POST("/tickets-status", handler.PostTicketsStatus)
	e.POST("/shows", handler.PostShows)
	e.POST("/book-tickets", handler.PostBookTickets)

	e.GET("/tickets", handler.GetTickets)

	return e
}
