package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
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

func (h *handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)

	w.WriteHeader(http.StatusCreated)
}

type getStatisticsResponse struct {
	Count int `json:"count"`
}

func (h *handler) handleGetStatistics(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)

	resp := getStatisticsResponse{Count: 0}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

// type errorResponse struct {
// 	Code    string `json:"code"`
// 	Message string `json:"message"`
// }

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
