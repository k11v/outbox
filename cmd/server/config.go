package main

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/k11v/squeak/internal/server"
)

const (
	modeDevelopment = "development"
	modeProduction  = "production"
)

type config struct {
	Mode   string        `env:"SQUEAK_MODE" envDefault:"development"`
	Server server.Config `envPrefix:"SQUEAK_SERVER_"`
}

func (c config) validate() error {
	if c.Mode != modeDevelopment && c.Mode != modeProduction {
		return fmt.Errorf("invalid mode: %s", c.Mode)
	}
	return c.Server.Validate()
}

func parseConfig(environ []string) (config, error) {
	cfg := config{}

	err := env.ParseWithOptions(&cfg, env.Options{
		Environment:     env.ToMap(environ),
		RequiredIfNoDef: true,
	})
	if err != nil {
		return config{}, err
	}

	err = cfg.validate()
	if err != nil {
		return config{}, err
	}

	return cfg, nil
}
