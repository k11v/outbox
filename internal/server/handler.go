package server

import (
	"encoding/json"
	"errors"
	"github.com/k11v/outbox/internal/outbox"
	"log/slog"
	"net/http"
)

type handler struct {
	log             *slog.Logger    // required
	messageProducer outbox.Producer // required
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

	headers := make([]outbox.Header, len(req.Headers))
	for i, header := range req.Headers {
		headers[i] = outbox.Header{
			Key:   header.Key,
			Value: []byte(header.Value),
		}
	}
	m := outbox.Message{
		Topic:   req.Topic,
		Key:     []byte(req.Key),
		Value:   []byte(req.Value),
		Headers: headers,
	}

	if err := h.messageProducer.Produce(r.Context(), m); err != nil {
		h.log.Error("failed to produce message", "error", err)
		if errors.Is(err, outbox.ErrUnknownTopicOrInternal) {
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
