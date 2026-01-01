package message

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"tickets/entities"
	"tickets/message/event"
)

func NewWatermillRouter(receiptsService event.ReceiptsService, spreadsheetsAPI event.SpreadsheetsAPI, rdb *redis.Client, watermillLogger watermill.LoggerAdapter) *message.Router {
	router := message.NewDefaultRouter(watermillLogger)

	handler := event.NewHandler(spreadsheetsAPI, receiptsService)

	useMiddlewares(router, watermillLogger)

	cqrsEventProcessor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return redisstream.NewSubscriber(
					redisstream.SubscriberConfig{
						Client:        rdb,
						ConsumerGroup: "ticket-service." + params.HandlerName,
						Consumer:      "consumer_1",
					},
					watermillLogger,
				)
			},
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return params.EventName, nil
			},
			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
			Logger: watermillLogger,
		},
	)
	if err != nil {
		panic(err)
	}

	cqrsEventProcessor.AddHandler(cqrs.NewEventHandler(
		"issue_receipt",
		func(ctx context.Context, event *entities.TicketBookingConfirmed) error {
			return handler.IssueReceipt(ctx, *event)
		},
	))
	cqrsEventProcessor.AddHandler(cqrs.NewEventHandler(
		"append_to_tracker",
		func(ctx context.Context, event *entities.TicketBookingConfirmed) error {
			return handler.AppendToTracker(ctx, *event)
		},
	))
	cqrsEventProcessor.AddHandler(cqrs.NewEventHandler(
		"cancel_ticket",
		func(ctx context.Context, event *entities.TicketBookingCanceled) error {
			return handler.CancelTicket(ctx, *event)
		},
	))

	return router
}
