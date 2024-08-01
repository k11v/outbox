package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/k11v/outbox/internal/kafkautil"
	"github.com/k11v/outbox/internal/postgresutil"
	"github.com/k11v/outbox/internal/worker"
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

	postgresPool, err := postgresutil.NewPool(ctx, log, cfg.Postgres, cfg.Development)
	if err != nil {
		return err
	}
	defer postgresPool.Close()

	w := worker.NewWorker(cfg.Worker, log, kafkaWriter, postgresPool)

	var done chan struct{}
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Info("interrupted")

		done <- struct{}{}
	}()

	log.Info(
		"starting worker",
		"development", cfg.Development,
	)
	w.Run(done)

	return nil
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
