package infra

import (
	"time"

	"github.com/IBM/sarama"
)

// EventProducer wraps the functionality of async producer.
type EventProducer struct {
	producer sarama.AsyncProducer
}

// NewEventProducer creates a new async producer.
func NewEventProducer(brokerList []string) (*EventProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal      // Only wait for leader to respond
	config.Producer.Flush.Frequency = 50 * time.Millisecond // Optional: Flush messages every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)

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
