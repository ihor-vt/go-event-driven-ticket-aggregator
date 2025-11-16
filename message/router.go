package message

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

func NewWatermillRouter(
	receiptsService ReceiptsService,
	spreadsheetsAPI SpreadsheetsAPI,
	rdb *redis.Client,
	watermillLogger watermill.LoggerAdapter,

) *message.Router {
	router := message.NewDefaultRouter(watermillLogger)

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-reciept",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-spreadsheet",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	router.AddConsumerHandler(
		"issue_receipt",
		"issue-receipt",
		issueReceiptSub,
		func(msg *message.Message) error {
			err := receiptsService.IssueReceipt(msg.Context(), string(msg.Payload))
			if err != nil {
				return fmt.Errorf("failed to issue receipt: %w", err)
			}

			return nil
		},
	)

	router.AddConsumerHandler(
		"append_to_tracker",
		"append-to-tracker",
		appendToTrackerSub,
		func(msg *message.Message) error {
			return spreadsheetsAPI.AppendRow(
				msg.Context(),
				"tickets-to-print",
				[]string{string(msg.Payload)},
			)
		},
	)

	return router
}
