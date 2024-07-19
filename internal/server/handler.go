package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type handler struct {
	log *slog.Logger
}

type getHealthResponse struct {
	Status string `json:"status"`
}

func (h *handler) handleGetHealth(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)

	resp := getHealthResponse{Status: "ok"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

type createMessageRequest struct {
	Content string `json:"content"`
}

type createMessageResponse struct {
	ID uuid.UUID `json:"id"`
}

func (h *handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)
	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := errorResponse{
			Code:    "invalid_request",
			Message: "invalid JSON",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			h.log.Error("failed to encode response", "error", err)
		}
		return
	}

	resp := createMessageResponse{ID: uuid.New()}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

type getStatisticsResponse struct {
	ProcessedMessages int `json:"processed_messages"`
}

func (h *handler) handleGetStatistics(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)

	resp := getStatisticsResponse{ProcessedMessages: 0}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
