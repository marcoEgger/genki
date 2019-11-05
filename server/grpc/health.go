package grpc

import (
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/lukasjarosch/genki/logger"
)

const (
	HealthUnknown    = 0
	HealthServing    = 1
	HealthNotServing = 2
)

// registerHealthServer will register a gRPC v1 Health server.
// The serving status will be set to NOT_SERVING.
func (srv *server) registerHealthServer() {
	srv.healthz = health.NewServer()
	srv.healthz.SetServingStatus(srv.opts.serviceName, HealthUnknown)
	grpc_health_v1.RegisterHealthServer(srv.grpc, srv.healthz)
	logger.Infof("registered gRPC health server for service '%s'", srv.opts.serviceName)
}

// setServingStatus of the gRPC server of the health server is enabled.
//
// 0 = UNKNOWN
// 1 = SERVING
// 2 = NOT_SERVING
func (srv *server) setServingStatus(status int32) {
	if srv.opts.HealthServerEnabled {
		srv.healthz.SetServingStatus(srv.opts.serviceName, grpc_health_v1.HealthCheckResponse_ServingStatus(status))
	}
}
