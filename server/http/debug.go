package http

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lukasjarosch/genki/logger"
)

// NewDebugServer is a convenience wrapper to quickly create a HTTP server, serving on port 3000.
// The server has the prometheus HTTP handler attached to '/metrics'.
func NewDebugServer(port string) Server {

	srv := NewServer(
		Name("debug"),
		Port(port),
		LoggingSkipEndpoints("/metrics"),
		LoggingSkipEndpoints("/health"),
	)

	srv.Handle("/metrics", promhttp.Handler())
	logger.Infof("prometheus metrics are exposed via 'debug' HTTP server on '/metrics'")

	return srv
}
