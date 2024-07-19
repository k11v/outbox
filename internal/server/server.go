package server

import (
	"crypto/tls"
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

func New(log *slog.Logger, cfg Config) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleGetHealth)

	subLogger := log.With("component", "server")
	subLogLogger := slog.NewLogLogger(subLogger.Handler(), slog.LevelError)

	return &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ErrorLog:          subLogLogger,
	}
}

func Listen(cfg Config) (net.Listener, error) {
	var err error
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	if !cfg.TLS.Enabled {
		return net.Listen("tcp", addr)
	}

	tlsCfg := &tls.Config{MinVersion: tls.VersionTLS13}
	tlsCfg.Certificates = make([]tls.Certificate, 1)
	tlsCfg.Certificates[0], err = tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.KeyFile)
	if err != nil {
		return nil, err
	}
	return tls.Listen("tcp", addr, tlsCfg)
}

func handleGetHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		panic(err)
	}
}
