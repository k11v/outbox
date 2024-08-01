package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (h *handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Store in Postgres.

	// TODO.

	// Send to Kafka.

	headers := make([]kafka.Header, len(req.Headers))
	for i, header := range req.Headers {
		headers[i] = kafka.Header{
			Key:   header.Key,
			Value: []byte(header.Value),
		}
	}
	m := kafka.Message{
		Topic:   req.Topic,
		Key:     []byte(req.Key),
		Value:   []byte(req.Value),
		Headers: headers,
	}

	if err := h.kafkaWriter.WriteMessages(r.Context(), m); err != nil {
		h.log.Error("failed to write message", "error", err)
		if errors.Is(err, kafka.UnknownTopicOrPartition) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type getStatisticsResponse struct {
	Count int `json:"count"`
}

func (h *handler) handleGetStatistics(w http.ResponseWriter, _ *http.Request) {
	resp := getStatisticsResponse{Count: 0}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to encode response", "error", err)
		return
	}
}
