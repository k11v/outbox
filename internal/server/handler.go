package server

import (
	"io"
	"log/slog"
	"net/http"
)

type handler struct {
	log *slog.Logger
}

func (h *handler) handleGetHealth(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (h *handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) handleGetStatistics(w http.ResponseWriter, r *http.Request) {
	defer closeWithLog(r.Body, h.log)
	w.WriteHeader(http.StatusOK)
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
