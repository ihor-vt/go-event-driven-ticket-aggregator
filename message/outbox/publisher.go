package outbox

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
)

func NewPublisherForDb(ctx context.Context, db *sqlx.Tx) (message.Publisher, error) {
	var publisher message.Publisher

	logger := watermill.NewSlogLogger(log.FromContext(ctx))

	publisher, err := watermillSQL.NewPublisher(
		db,
		watermillSQL.PublisherConfig{
			SchemaAdapter: watermillSQL.DefaultPostgreSQLSchema{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create outbox publisher: %w", err)
	}
	publisher = log.CorrelationPublisherDecorator{publisher}

	publisher = forwarder.NewPublisher(publisher, forwarder.PublisherConfig{
		ForwarderTopic: outboxTopic,
	})
	publisher = log.CorrelationPublisherDecorator{publisher}

	return publisher, nil
}
