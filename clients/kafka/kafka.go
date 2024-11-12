package kafka

import (
	"context"

	"github.com/erikqwerty/erik-platform/clients/kafka/consumer"
)

// Consumer - интерфейс для обработки чтения сообщений из Kafka
type Consumer interface {
	Consume(ctx context.Context, topicName string, handler consumer.Handler) (err error)
	Close() error
}

// Producer - интерфейс для отправки сообщений в Kafka.
type Producer interface {
	SendMessage(topic string, value string) (partition int32, offset int64, err error)
	Close() error
}
