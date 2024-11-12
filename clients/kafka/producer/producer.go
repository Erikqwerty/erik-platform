package producer

import (
	"github.com/IBM/sarama"

	"github.com/erikqwerty/erik-platform/clients/kafka"
)

// syncProducer - реализация интерфейса Producer с использованием sarama.SyncProducer.
type syncProducer struct {
	producer sarama.SyncProducer
}

// NewProducer - конструктор, создающий новый синхронный продюсер.
func NewProducer(brokers []string) (kafka.Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &syncProducer{producer: producer}, nil
}

// SendMessage - отправляет сообщение в указанный топик. возвращает partition, offset, err
func (p *syncProducer) SendMessage(topic string, value string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	return partition, offset, err
}

// Close - закрывает соединение продюсера.
func (p *syncProducer) Close() error {
	return p.producer.Close()
}
