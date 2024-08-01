package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k11v/outbox/internal/outbox"
	"github.com/segmentio/kafka-go"

	"github.com/k11v/outbox/internal/kafkautil"
	"github.com/k11v/outbox/internal/postgresutil"
)

func main() {
	if err := run(os.Stdout, os.Environ()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(stdout io.Writer, environ []string) error {
	cfg, err := parseConfig(environ)
	if err != nil {
		return err
	}
	log := newLogger(stdout, cfg.Development)

	ctx := context.Background()

	kafkaWriter := kafkautil.NewWriter(cfg.Kafka)
	defer closeWithLog(kafkaWriter, log)

	postgresPool, err := postgresutil.NewPool(ctx, log, cfg.Postgres, cfg.Development)
	if err != nil {
		return err
	}
	defer postgresPool.Close()

	w := struct {
		log          *slog.Logger
		kafkaWriter  *kafka.Writer
		postgresPool *pgxpool.Pool
	}{
		log:          log,
		kafkaWriter:  kafkaWriter,
		postgresPool: postgresPool,
	}

	ticker := time.NewTicker(cfg.Worker.interval())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	log.Info(
		"starting worker",
		"development", cfg.Development,
	)

	for {
		var workErr error
		workCtx := context.Background()

		// Retrieve from Postgres.

		type row struct {
			ID      uuid.UUID `json:"id"`
			Topic   string    `json:"topic"`
			Key     string    `json:"key"`
			Value   string    `json:"value"`
			Headers []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}
		}

		result, workErr := postgresPool.Query(
			workCtx,
			`
			SELECT id, topic, key, value, headers::jsonb
			FROM outbox_messages
			WHERE status = $1
			ORDER BY created_at, id
			LIMIT $2
		`,
			outbox.StatusUndelivered,
			cfg.Worker.batchSize(),
		)
		if workErr != nil {
			w.log.Error("failed to query outbox_messages", "error", err)
			continue
		}

		rows, workErr := pgx.CollectRows(result, pgx.RowToStructByName[row])
		if workErr != nil {
			w.log.Error("failed to collect rows", "error", err)
			continue
		}

		if len(rows) == 0 {
			w.log.Debug("no messages to send")
			select {
			case <-ticker.C:
			case <-done:
				return nil
			}
			continue
		}

		// Send to Kafka.

		w.log.Info("sending messages", "count", len(rows))

		var messages []kafka.Message
		for _, mr := range rows {
			headers := make([]kafka.Header, len(mr.Headers))
			for i, header := range mr.Headers {
				headers[i] = kafka.Header{
					Key:   header.Key,
					Value: []byte(header.Value),
				}
			}
			m := kafka.Message{
				Topic:   mr.Topic,
				Key:     []byte(mr.Key),
				Value:   []byte(mr.Value),
				Headers: headers,
			}
			messages = append(messages, m)
		}

		if err = w.kafkaWriter.WriteMessages(workCtx, messages...); err != nil {
			w.log.Error("failed to write message", "error", err)
			continue
		}

		// Update in Postgres.

		var ids []uuid.UUID
		for _, mr := range rows {
			ids = append(ids, mr.ID)
		}

		_, err = w.postgresPool.Exec(
			workCtx,
			`UPDATE outbox_messages SET status = $1 WHERE id = ANY($2)`,
			outbox.StatusDelivered,
			ids,
		)
		if err != nil {
			w.log.Error("failed to update outbox_messages", "error", err)
			continue
		}

		select {
		case <-ticker.C:
		case <-done:
			return nil
		}
	}
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
