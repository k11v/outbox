package server

import (
	"crypto/tls"
	"github.com/k11v/outbox/internal/outbox"
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

// New returns a new HTTP server.
// It should be started with a listener returned by Listen.
func New(cfg Config, log *slog.Logger, messageProducer outbox.Producer) *http.Server {
	mux := http.NewServeMux()

	h := &handler{log: log, messageProducer: messageProducer}
	mux.HandleFunc("GET /health", h.handleGetHealth)
	mux.HandleFunc("POST /messages", h.handleCreateMessage)
	mux.HandleFunc("GET /statistics", h.handleGetStatistics)

	subLogger := log.With("component", "server")
	subLogLogger := slog.NewLogLogger(subLogger.Handler(), slog.LevelError)

	return &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ErrorLog:          subLogLogger,
	}
}

// Listen listens on the TCP network address addr and returns a net.Listener.
// If TLS is enabled, it listens for TLS connections.
func Listen(cfg Config) (net.Listener, error) {
	var err error
	addr := net.JoinHostPort(cfg.host(), strconv.Itoa(cfg.port()))

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
