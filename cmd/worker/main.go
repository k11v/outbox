package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/k11v/outbox/internal/worker"

	"github.com/k11v/outbox/internal/kafkautil"
	"github.com/k11v/outbox/internal/postgresutil"
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

	ticker := time.NewTicker(cfg.Worker.IntervalReal())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	log.Info(
		"starting worker",
		"development", cfg.Development,
	)

	for {
		workCtx := context.Background()

		sendCount, sendErr := w.SendMessages(workCtx)
		if sendErr != nil {
			log.Error("failed to send messages", "error", sendErr)
		} else {
			if sendCount > 0 {
				log.Info("sent messages", "count", sendCount)
			} else {
				log.Debug("sent no messages")
			}
		}

		select {
		case <-ticker.C:
		case <-done:
			return nil
		}
	}
}

func closeWithLog(c io.Closer, log *slog.Logger) {
	if err := c.Close(); err != nil {
		log.Error("failed to close", "error", err)
	}
}
