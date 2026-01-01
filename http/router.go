package http

import (
	"net/http"
	"tickets/message"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
	"github.com/lithammer/shortuuid/v3"
)

func NewHttpRouter(eventBus *cqrs.EventBus) *echo.Echo {
	e := libHttp.NewEcho()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			correlationId := req.Header.Get("Correlation-ID")

			if correlationId == "" {
				correlationId = shortuuid.New()
			}

			ctx := message.ContextWithCorrelationId(req.Context(), correlationId)

			c.SetRequest(req.WithContext(ctx))

			return next(c)
		}
	})

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eventBus: *eventBus,
	}

	e.POST("/tickets-status", handler.PostTicketsStatus)

	return e
}
