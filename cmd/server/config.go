package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/k11v/squeak/internal/server"
)

// config holds the application configuration.
// The zero value is a valid configuration.
type config struct {
	Development bool          `env:"SQUEAK_DEVELOPMENT"`
	Server      server.Config `envPrefix:"SQUEAK_SERVER_"`
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
