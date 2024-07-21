package message

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
)

var (
	ErrUnknownTopic = errors.New("unknown topic")
)

type Producer interface {
	Produce(ctx context.Context, messages []Message) error
}

type KafkaProducer struct {
	Writer *kafka.Writer // required, its Topic must be unset
}

func (p *KafkaProducer) Produce(ctx context.Context, messages []Message) error {
	kafkaMessages := make([]kafka.Message, len(messages))
	for i, m := range messages {
		kafkaMessages[i] = kafka.Message{
			Topic: m.Topic,
			Key:   m.Key,
			Value: m.Value,
		}
	}

	if err := p.Writer.WriteMessages(ctx, kafkaMessages...); err != nil {
		if errors.Is(err, kafka.UnknownTopicOrPartition) {
			return errors.Join(ErrUnknownTopic, err)
		}
		return err
	}

	return nil
}
