package main

import (
	"fmt"
	"io"
	"log/slog"
)

func newLogger(w io.Writer, mode string) (*slog.Logger, error) {
	var handler slog.Handler
	switch mode {
	case modeDevelopment:
		opts := &slog.HandlerOptions{Level: slog.LevelDebug}
		handler = slog.NewTextHandler(w, opts)
	case modeProduction:
		opts := &slog.HandlerOptions{Level: slog.LevelInfo}
		handler = slog.NewJSONHandler(w, opts)
	default:
		return nil, fmt.Errorf("invalid mode: %s", mode)
	}

	return slog.New(handler), nil
}
