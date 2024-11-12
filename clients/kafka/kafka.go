package kafka

import (
	"context"

	"github.com/erikqwerty/erik-platform/clients/kafka/consumer"
)

type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}
