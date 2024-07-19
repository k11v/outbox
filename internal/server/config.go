package server

import (
	"time"
)

type Config struct {
	Host              string        `env:"HOST"` // default: "127.0.0.1"
	Port              int           `env:"PORT"` // default: 8080
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT"`
	TLS               TLSConfig     `envPrefix:"TLS_"`
}

type TLSConfig struct {
	Enabled  bool   `env:"ENABLED"`
	CertFile string `env:"CERT_FILE"`
	KeyFile  string `env:"KEY_FILE"`
}

func (c Config) host() string {
	h := c.Host
	if h == "" {
		h = "127.0.0.1"
	}
	return h
}

func (c Config) port() int {
	p := c.Port
	if p == 0 {
		p = 8080
	}
	return p
}
