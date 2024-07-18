package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/k11v/squeak/internal/server"
)

func main() {
	if err := run(os.Stdout, os.Getenv); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(stdout io.Writer, getenv func(string) string) error {
	cfg, err := parseConfig(getenv)
	if err != nil {
		return err
	}

	log, err := newLogger(stdout, cfg.Mode)
	if err != nil {
		return err
	}

	srv := server.New(log, cfg.Server)

	lst, err := server.NewListener(cfg.Server)
	if err != nil {
		return err
	}

	log.Info(
		"starting server",
		"mode", cfg.Mode,
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"tls", cfg.Server.TLS.Enabled,
	)
	if err = srv.Serve(lst); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
