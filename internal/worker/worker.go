package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k11v/outbox/internal/outbox"
	"github.com/segmentio/kafka-go"
)

type Worker struct {
	cfg          Config
	log          *slog.Logger
	kafkaWriter  *kafka.Writer
	postgresPool *pgxpool.Pool
}

func NewWorker(cfg Config, log *slog.Logger, kafkaWriter *kafka.Writer, postgresPool *pgxpool.Pool) *Worker {
	return &Worker{
		cfg:          cfg,
		kafkaWriter:  kafkaWriter,
		log:          log,
		postgresPool: postgresPool,
	}
}

func (w *Worker) Run(done <-chan struct{}) {
	ticker := time.NewTicker(w.cfg.interval())

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), w.cfg.timeout())
			defer cancel()

			count, err := w.sendMessages(ctx)
			if err != nil {
				w.log.Error("failed to send messages", "error", err)
				return
			}

			if count > 0 {
				w.log.Info("sent messages", "count", count)
			} else {
				w.log.Debug("sent no messages")
			}
		}()

		select {
		case <-ticker.C:
		case <-done:
			break
		}
	}
}

func (w *Worker) sendMessages(ctx context.Context) (int, error) {
	// Get undelivered messages.

	result, err := w.postgresPool.Query(
		ctx,
		`
			SELECT id, topic, key, value, headers::jsonb
			FROM outbox_messages
			WHERE status = $1
			ORDER BY created_at, id
			LIMIT $2
		`,
		outbox.StatusUndelivered,
		w.cfg.batchSize(),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to query outbox_messages: %w", err)
	}

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
	rows, err := pgx.CollectRows(result, pgx.RowToStructByName[row])
	if err != nil {
		return 0, fmt.Errorf("failed to collect rows: %w", err)
	}

	if len(rows) == 0 {
		return 0, nil
	}

	// Send messages.

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

	if err = w.kafkaWriter.WriteMessages(ctx, messages...); err != nil {
		return 0, fmt.Errorf("failed to write messages: %w", err)
	}

	// Update status of messages.

	var ids []uuid.UUID
	for _, mr := range rows {
		ids = append(ids, mr.ID)
	}

	_, err = w.postgresPool.Exec(
		ctx,
		`UPDATE outbox_messages SET status = $1 WHERE id = ANY($2)`,
		outbox.StatusDelivered,
		ids,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to update outbox_messages: %w", err)
	}

	return len(messages), nil
}
