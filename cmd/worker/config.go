package main

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/k11v/outbox/internal/kafkautil"
	"github.com/k11v/outbox/internal/postgresutil"
)

// config holds the application configuration.
type config struct {
	Development bool                `env:"OUTBOX_DEVELOPMENT"`
	Kafka       kafkautil.Config    `envPrefix:"OUTBOX_KAFKA_"`
	Postgres    postgresutil.Config `envPrefix:"OUTBOX_POSTGRES_"`
	Worker      workerConfig        `envPrefix:"OUTBOX_WORKER_"`
}

type workerConfig struct {
	BatchSize int           `env:"BATCH_SIZE"` // default: 100
	Interval  time.Duration `env:"INTERVAL"`   // default: 1s
}

func (c workerConfig) batchSize() int {
	s := c.BatchSize
	if s == 0 {
		s = 100
	}
	return s
}

func (c workerConfig) interval() time.Duration {
	i := c.Interval
	if i == 0 {
		i = time.Second
	}
	return i
}

// parseConfig parses the application configuration from the environment variables.
func parseConfig(environ []string) (config, error) {
	cfg := config{}

	err := env.ParseWithOptions(&cfg, env.Options{
		Environment: env.ToMap(environ),
	})
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}
