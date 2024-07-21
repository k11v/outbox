package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	log := slog.Default()

	topic := "example"
	addr := "kafka:9092"

	if err := createTopic(topic, addr, log); err != nil {
		return err
	}

	return nil
}

func createTopic(topic, addr string, log *slog.Logger) error {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer closeWithLog(conn, log)

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer closeWithLog(controllerConn, log)

	topicCfg := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	if err = controllerConn.CreateTopics(topicCfg); err != nil {
		return err
	}

	return nil
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
