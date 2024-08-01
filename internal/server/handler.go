package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/k11v/outbox/internal/outbox"
	"github.com/segmentio/kafka-go"
)

type handler struct {
	kafkaWriter  *kafka.Writer // required, its Topic must be unset
	log          *slog.Logger  // required
	postgresPool *pgxpool.Pool // required
}

type getHealthResponse struct {
	Status string `json:"status"`
}

func (h *handler) handleGetHealth(w http.ResponseWriter, _ *http.Request) {
	resp := getHealthResponse{Status: "ok"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

type createMessageRequest struct {
	Topic   string                       `json:"topic"`
	Key     string                       `json:"key"`
	Value   string                       `json:"value"`
	Headers []createMessageHeaderRequest `json:"headers"`
}

type createMessageHeaderRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (r *createMessageRequest) validate() error {
	if r.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if r.Key == "" {
		return fmt.Errorf("key is required")
	}
	if r.Value == "" {
		return fmt.Errorf("value is required")
	}
	for i, header := range r.Headers {
		if header.Key == "" {
			return fmt.Errorf("header key is required at index %d", i)
		}
	}
	return nil
}

func (h *handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("failed to decode request: %v", err)))
		return
	}
	if err := req.validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf("invalid request: %v", err)))
		return
	}

	if err := h.createMessage(r.Context(), req); err != nil {
		h.log.Error("failed to create message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) createMessage(ctx context.Context, req createMessageRequest) error {
	tx, err := h.postgresPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func(tx pgx.Tx) {
		_ = tx.Rollback(ctx)
	}(tx)

	// Insert into message_infos just to have a need for the transactional outbox pattern.

	_, err = tx.Exec(
		ctx,
		`INSERT INTO message_infos (value_length) VALUES ($1)`,
		len(req.Value),
	)
	if err != nil {
		return fmt.Errorf("failed to insert into message_infos: %w", err)
	}

	// Insert outbox_messages to have a message to send to Kafka by the worker.

	headersJSON, err := json.Marshal(req.Headers)
	if err != nil {
		return fmt.Errorf("failed to marshal headers: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO outbox_messages (status, topic, key, value, headers) VALUES ($1, $2, $3, $4, $5)`,
		outbox.StatusUndelivered,
		req.Topic,
		req.Key,
		req.Value,
		headersJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to insert into outbox_messages: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

type getStatisticsResponse struct {
	CountInMessageInfos              int `json:"count_in_message_infos"`
	UndeliveredCountInOutboxMessages int `json:"undelivered_count_in_outbox_messages"`
	DeliveredCountInOutboxMessages   int `json:"delivered_count_in_outbox_messages"`
}

func (h *handler) handleGetStatistics(w http.ResponseWriter, r *http.Request) {
	resp, err := h.getStatistics(r.Context())
	if err != nil {
		h.log.Error("failed to get statistics", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to encode response", "error", err)
		return
	}
}

func (h *handler) getStatistics(ctx context.Context) (*getStatisticsResponse, error) {
	result, err := h.postgresPool.Query(
		ctx,
		`
			SELECT
				(SELECT COUNT(*) FROM message_infos) AS count_in_message_infos,
				(SELECT COUNT(*) FROM outbox_messages WHERE status = $1) AS undelivered_count_in_outbox_messages,
				(SELECT COUNT(*) FROM outbox_messages WHERE status = $2) AS delivered_count_in_outbox_messages
		`,
		outbox.StatusUndelivered,
		outbox.StatusDelivered,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query statistics: %w", err)
	}

	type row struct {
		CountInMessageInfos              int `db:"count_in_message_infos"`
		UndeliveredCountInOutboxMessages int `db:"undelivered_count_in_outbox_messages"`
		DeliveredCountInOutboxMessages   int `db:"delivered_count_in_outbox_messages"`
	}
	r, err := pgx.CollectExactlyOneRow(result, pgx.RowToStructByName[row])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &getStatisticsResponse{
		CountInMessageInfos:              r.CountInMessageInfos,
		UndeliveredCountInOutboxMessages: r.UndeliveredCountInOutboxMessages,
		DeliveredCountInOutboxMessages:   r.DeliveredCountInOutboxMessages,
	}, nil
}
