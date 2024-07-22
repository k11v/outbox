package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/k11v/outbox/internal/kafkautil"
)

// config holds the application configuration.
type config struct {
	Kafka kafkautil.Config `envPrefix:"OUTBOX_KAFKA_"`
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
