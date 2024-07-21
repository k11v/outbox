package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/k11v/squeak/internal/message"
)

type handler struct {
	log             *slog.Logger     // required
	messageProducer message.Producer // required
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
	Topic string `json:"topic"`
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

	m := message.Message{
		Topic: req.Topic,
		Key:   []byte(req.Key),
		Value: []byte(req.Value),
	}

	if err := h.messageProducer.Produce(r.Context(), []message.Message{m}); err != nil {
		// FIXME: Missing topic could be handled differently.
		h.log.Error("failed to produce message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
	}
}
