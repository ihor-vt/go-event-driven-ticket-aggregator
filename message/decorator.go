package message

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lithammer/shortuuid/v3"
)

type ctxKey int

const (
	correlationIdKey ctxKey = iota
)

func ContextWithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, correlationIdKey, correlationId)
}

func CorrelationIdFromContext(ctx context.Context) string {
	v, ok := ctx.Value(correlationIdKey).(string)
	if ok {
		return v
	}

	return "gen_" + shortuuid.New()
}

type CorrelationPublisherDecorator struct {
	message.Publisher
}

func (c CorrelationPublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	for i := range messages {
		correlationId := CorrelationIdFromContext(messages[i].Context())
		messages[i].Metadata.Set("correlation_id", correlationId)
	}

	return c.Publisher.Publish(topic, messages...)
}
