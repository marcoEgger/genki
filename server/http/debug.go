package http

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/lukasjarosch/genki/logger"
)

// NewDebugServer is a convenience wrapper to quickly create a HTTP server, serving on port 3000.
// The server has the prometheus HTTP handler attached to '/metrics'.
func NewDebugServer() Server {
	srv := NewServer(
		Name("debug"),
		Port("3000"),
	)

	srv.Handle("/metrics", promhttp.Handler())
	logger.Debugf("prometheus metrics are exposed via 'debug' on '/metrics'")

	return srv
}