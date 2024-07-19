package server

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealth(t *testing.T) {
	t.Run("Returns OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()

		srv := New(slog.Default(), Config{})
		srv.Handler.ServeHTTP(rec, req)

		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := rec.Body.String(), `{"status":"ok"}`; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
