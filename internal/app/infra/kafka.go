package infra

import (
	"github.com/IBM/sarama"
)

// EventProducer wraps the functionality of async producer.
type EventProducer struct {
	producer sarama.AsyncProducer
}

// NewEventProducer creates a new async producer.
func NewEventProducer(brokerList []string) (*EventProducer, error) {
	producer, err := sarama.NewAsyncProducer(brokerList, nil)

	return &EventProducer{producer: producer}, err
}

// SendMessage sends a message to the broker.
func (p *EventProducer) SendMessage(topic string, payload []byte) error {
	kafkaMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	}

	p.producer.Input() <- kafkaMessage

	select {
	case <-p.producer.Successes():
		return nil
	case err := <-p.producer.Errors():
		return err
	}
}
