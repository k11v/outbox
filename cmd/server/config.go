package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/k11v/squeak/internal/server"
)

type config struct {
	Development bool          `env:"DEVELOPMENT"`
	Server      server.Config `envPrefix:"SQUEAK_SERVER_"`
}

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
