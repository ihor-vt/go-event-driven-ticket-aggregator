package message

import (
	"log/slog"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func useMiddlewares(router *message.Router, watermilLogger watermill.LoggerAdapter) {
	router.AddMiddleware(middleware.Recoverer)

	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			correlationID := msg.Metadata.Get("correlation_id")
			if correlationID == "" {
				correlationID = msg.UUID
			}

			ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
			msg.SetContext(ctx)
			return next(msg)
		}

	})

	router.AddMiddleware(func(next message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			logger := log.FromContext(msg.Context()).With(
				slog.String("message_id", msg.UUID),
				slog.String("payload", string(msg.Payload)),
				slog.Any("metadata", msg.Metadata),
				slog.String("handler", message.HandlerNameFromCtx(msg.Context())),
			)

			logger.Info("Handling a message")
			msgs, err := next(msg)
			if err != nil {
				logger.Error("Error while handling a message", slog.String("error", err.Error()))
			}

			return msgs, err
		}
	})

	router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Microsecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermilLogger,
	}.Middleware)
}
