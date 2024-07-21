package message

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(ctx context.Context, messages []Message) error
}

type KafkaProducer struct {
	Writer *kafka.Writer // required
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

	err := p.Writer.WriteMessages(ctx, kafkaMessages...)
	if err != nil {
		return err
	}

	return nil
}

func NewKafkaWriter(addr string) *kafka.Writer {
	addrs := strings.Split(addr, ",")
	return &kafka.Writer{
		Addr:         kafka.TCP(addrs...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
}
