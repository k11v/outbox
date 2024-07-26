package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetHealth(t *testing.T) {
	t.Run("Returns OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()

		// FIXME: Config is empty, kafkaWriter and postgresPool are nil.
		srv := New(Config{}, slog.Default(), nil, nil)
		srv.Handler.ServeHTTP(rec, req)

		if got, want := rec.Code, http.StatusOK; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := rec.Body.String(), `{"status":"ok"}`; !equalJSON(got, want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func equalJSON(x, y string) bool {
	var mx, my any
	if err := json.Unmarshal([]byte(x), &mx); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(y), &my); err != nil {
		return false
	}
	return reflect.DeepEqual(mx, my)
}
