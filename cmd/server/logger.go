package main

import (
	"fmt"
	"io"
	"log/slog"
)

func newLogger(w io.Writer, mode string) *slog.Logger {
	var handler slog.Handler
	switch mode {
	case modeDevelopment:
		opts := &slog.HandlerOptions{Level: slog.LevelDebug}
		handler = slog.NewTextHandler(w, opts)
	case modeProduction:
		opts := &slog.HandlerOptions{Level: slog.LevelInfo}
		handler = slog.NewJSONHandler(w, opts)
	default:
		panic(fmt.Sprintf("unknown mode: %s", mode))
	}

	return slog.New(handler)
}
