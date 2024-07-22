package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/k11v/squeak/internal/kafkautil"
	"github.com/k11v/squeak/internal/message"
	"github.com/k11v/squeak/internal/postgresutil"
	"github.com/k11v/squeak/internal/server"
)

func main() {
	if err := run(os.Stdout, os.Environ()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(stdout io.Writer, environ []string) error {
	cfg, err := parseConfig(environ)
	if err != nil {
		return err
	}
	log := newLogger(stdout, cfg.Development)

	ctx := context.Background()

	kafkaWriter := kafkautil.NewWriter(cfg.Kafka)
	defer closeWithLog(kafkaWriter, log)

	messageProducer := &message.KafkaProducer{Writer: kafkaWriter}

	postgresPool, err := postgresutil.NewPool(ctx, log, cfg.Postgres, cfg.Development)
	if err != nil {
		return err
	}
	defer postgresPool.Close()

	srv := server.New(cfg.Server, log, messageProducer)
	lst, err := server.Listen(cfg.Server)
	if err != nil {
		return err
	}

	log.Info(
		"starting server",
		"addr", lst.Addr().String(),
		"development", cfg.Development,
		"tls", cfg.Server.TLS.Enabled,
	)
	if err = srv.Serve(lst); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
