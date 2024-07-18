package server

import "time"

type Config struct {
	Host              string
	Port              int
	ReadHeaderTimeout time.Duration
	TLS               TLSConfig
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}
