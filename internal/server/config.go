package server

import (
	"fmt"
	"time"
)

type Config struct {
	Host              string        `env:"HOST" envDefault:"127.0.0.1"`
	Port              int           `env:"PORT" envDefault:"8000"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"1s"`
	TLS               TLSConfig     `envPrefix:"TLS_"`
}

func (c Config) Validate() error {
	if !(0 < c.Port && c.Port < 65536) {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	return c.TLS.Validate()
}

type TLSConfig struct {
	Enabled  bool   `env:"ENABLED" envDefault:"false"`
	CertFile string `env:"CERT_FILE" envDefault:""`
	KeyFile  string `env:"KEY_FILE" envDefault:""`
}

func (c TLSConfig) Validate() error {
	if c.Enabled && (c.CertFile == "" || c.KeyFile == "") {
		return fmt.Errorf("missing cert or key file with TLS enabled")
	}
	return nil
}
