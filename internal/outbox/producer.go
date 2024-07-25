package outbox

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
)

var (
	ErrUnknownTopicOrInternal = errors.New("unknown topic or internal")
)

type Producer interface {
	Produce(ctx context.Context, messages ...Message) error
}

type KafkaProducer struct {
	Writer *kafka.Writer // required, its Topic must be unset
}

func (p *KafkaProducer) Produce(ctx context.Context, messages ...Message) error {
	kafkaMessages := make([]kafka.Message, len(messages))
	for i, m := range messages {
		headers := make([]kafka.Header, len(m.Headers))
		for j, h := range m.Headers {
			headers[j] = kafka.Header{
				Key:   h.Key,
				Value: h.Value,
			}
		}

		kafkaMessages[i] = kafka.Message{
			Topic:   m.Topic,
			Key:     m.Key,
			Value:   m.Value,
			Headers: headers,
		}
	}

	if err := p.Writer.WriteMessages(ctx, kafkaMessages...); err != nil {
		if errors.Is(err, kafka.UnknownTopicOrPartition) {
			return errors.Join(ErrUnknownTopicOrInternal, err)
		}
		return err
	}

	return nil
}
