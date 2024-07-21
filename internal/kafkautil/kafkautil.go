package kafkautil

import (
	"strings"

	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(addr string) *kafka.Writer {
	addrs := strings.Split(addr, ",")
	return &kafka.Writer{
		Addr:         kafka.TCP(addrs...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
}
