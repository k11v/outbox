package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

var topicCfgs = []kafka.TopicConfig{
	{
		Topic:             "example",
		NumPartitions:     1,
		ReplicationFactor: 1,
	},
	{
		Topic:             "example2",
		NumPartitions:     1,
		ReplicationFactor: 1,
	},
	{
		Topic:             "example3",
		NumPartitions:     1,
		ReplicationFactor: 1,
	},
}

func main() {
	if err := run(os.Environ()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(environ []string) error {
	cfg, err := parseConfig(environ)
	if err != nil {
		return err
	}
	log := slog.Default()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err = createTopics(ctx, log, cfg.Kafka.Brokers, topicCfgs...); err != nil {
		return err
	}

	return nil
}

func createTopics(ctx context.Context, log *slog.Logger, brokers []string, topicCfgs ...kafka.TopicConfig) error {
	// Connect to the first broker that is available.

	var conn *kafka.Conn
	var err error
	for _, broker := range brokers {
		conn, err = kafka.DialContext(ctx, "tcp", broker)
		if err == nil {
			break
		}
	}
	if err != nil {
		return errors.Join(errors.New("failed to connect to any broker"), err)
	}
	defer closeWithLog(conn, log)

	// Get the controller broker.

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	// Connect to the controller broker.

	controllerConn, err := kafka.DialContext(ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return errors.Join(errors.New("failed to connect to controller"), err)
	}
	defer closeWithLog(controllerConn, log)

	// Create the topics.

	if err = controllerConn.CreateTopics(topicCfgs...); err != nil {
		return err
	}

	return nil
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
