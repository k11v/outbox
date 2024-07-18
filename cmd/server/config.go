package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/k11v/squeak/internal/server"
)

const (
	modeKey                    = "SQUEAK_MODE"
	serverHostKey              = "SQUEAK_SERVER_HOST"
	serverPortKey              = "SQUEAK_SERVER_PORT"
	serverReadHeaderTimeoutKey = "SQUEAK_SERVER_READ_HEADER_TIMEOUT"
	serverTLSEnabledKey        = "SQUEAK_SERVER_TLS_ENABLED"
	serverTLSCertFileKey       = "SQUEAK_SERVER_TLS_CERT_FILE"
	serverTLSKeyFileKey        = "SQUEAK_SERVER_TLS_KEY_FILE"
)

const (
	defaultMode                    = modeDevelopment
	defaultServerHost              = "127.0.0.1"
	defaultServerPort              = 8000
	defaultServerReadHeaderTimeout = 1 * time.Second
	defaultServerTLSEnabled        = false
)

const (
	modeDevelopment = "development"
	modeProduction  = "production"
)

type config struct {
	Mode   string
	Server server.Config
}

func parseConfig(getenv func(string) string) (config, error) {
	var err error

	mode := defaultMode
	if getenv(modeKey) != "" {
		mode = getenv(modeKey)
		switch mode {
		case modeDevelopment, modeProduction:
		default:
			return config{}, errors.New("invalid mode")
		}
	}

	serverCfg, err := parseServerConfig(getenv)
	if err != nil {
		return config{}, err
	}

	return config{
		Mode:   mode,
		Server: serverCfg,
	}, nil
}

func parseServerConfig(getenv func(string) string) (server.Config, error) {
	var err error

	host := defaultServerHost
	if getenv(serverHostKey) != "" {
		host = getenv(serverHostKey)
	}

	port := defaultServerPort
	if getenv(serverPortKey) != "" {
		port, err = strconv.Atoi(getenv(serverPortKey))
		if err != nil {
			return server.Config{}, errors.Join(errors.New("invalid server port"), err)
		}
	}

	readHeaderTimeout := defaultServerReadHeaderTimeout
	if getenv(serverReadHeaderTimeoutKey) != "" {
		readHeaderTimeout, err = time.ParseDuration(getenv(serverReadHeaderTimeoutKey))
		if err != nil {
			return server.Config{}, errors.Join(errors.New("invalid server read header timeout"), err)
		}
	}

	tlsCfg, err := parseServerTLSConfig(getenv)
	if err != nil {
		return server.Config{}, err
	}

	return server.Config{
		Host:              host,
		Port:              port,
		ReadHeaderTimeout: readHeaderTimeout,
		TLS:               tlsCfg,
	}, nil
}

func parseServerTLSConfig(getenv func(string) string) (server.TLSConfig, error) {
	var err error

	enabled := defaultServerTLSEnabled
	if getenv(serverTLSEnabledKey) != "" {
		enabled, err = strconv.ParseBool(getenv(serverTLSEnabledKey))
		if err != nil {
			return server.TLSConfig{}, errors.Join(errors.New("invalid server tls enabled"), err)
		}
	}

	certFile := ""
	if getenv(serverTLSCertFileKey) != "" {
		certFile = getenv(serverTLSCertFileKey)
	}

	keyFile := ""
	if getenv(serverTLSKeyFileKey) != "" {
		keyFile = getenv(serverTLSKeyFileKey)
	}

	if enabled && (certFile == "" || keyFile == "") {
		return server.TLSConfig{}, errors.New("missing server tls cert file or key file")
	}

	return server.TLSConfig{
		Enabled:  enabled,
		CertFile: certFile,
		KeyFile:  keyFile,
	}, nil
}
