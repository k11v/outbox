package kafkautil

import (
	"github.com/segmentio/kafka-go"
)

// NewKafkaWriter creates a new kafka.Writer.
func NewKafkaWriter(cfg Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
}
