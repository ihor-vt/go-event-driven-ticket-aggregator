package message

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"

	"tickets/message/command"
	"tickets/message/event"
	"tickets/message/outbox"
)

func NewWatermillRouter(
	postgresSubscriber message.Subscriber,
	publisher message.Publisher,
	eventProcessorConfig cqrs.EventProcessorConfig,
	eventHandler event.Handler,
	commandProcessorConfig cqrs.CommandProcessorConfig,
	commandsHandler command.Handler,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router := message.NewDefaultRouter(watermillLogger)

	useMiddlewares(router, watermillLogger)

	outbox.AddForwarderHandler(postgresSubscriber, publisher, router, watermillLogger)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"BookPlaceInDeadNation",
			eventHandler.BookPlaceInDeadNation,
		),
		cqrs.NewEventHandler(
			"AppendToTracker",
			eventHandler.AppendToTracker,
		),
		cqrs.NewEventHandler(
			"TicketRefundToSheet",
			eventHandler.TicketRefundToSheet,
		),
		cqrs.NewEventHandler(
			"IssueReceipt",
			eventHandler.IssueReceipt,
		),
		cqrs.NewEventHandler(
			"PrintTicketHandler",
			eventHandler.PrintTicket,
		),
		cqrs.NewEventHandler(
			"StoreTickets",
			eventHandler.StoreTickets,
		),
		cqrs.NewEventHandler(
			"RemoveCanceledTicket",
			eventHandler.RemoveCanceledTicket,
		),
	)

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(router, commandProcessorConfig)
	if err != nil {
		panic(err)
	}

	commandProcessor.AddHandlers(
		cqrs.NewCommandHandler(
			"TicketRefund",
			commandsHandler.RefundTicket,
		),
	)

	return router
}
