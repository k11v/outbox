package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/k11v/outbox/internal/kafkautil"
	"github.com/k11v/outbox/internal/postgresutil"
	"github.com/k11v/outbox/internal/server"
)

// config holds the application configuration.
type config struct {
	Development bool                `env:"OUTBOX_DEVELOPMENT"`
	Kafka       kafkautil.Config    `envPrefix:"OUTBOX_KAFKA_"`
	Postgres    postgresutil.Config `envPrefix:"OUTBOX_POSTGRES_"`
	Server      server.Config       `envPrefix:"OUTBOX_SERVER_"`
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
