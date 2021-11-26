package grpc

import (
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/marcoEgger/genki/logger"
)

const (
	HealthUnknown    = 0
	HealthServing    = 1
	HealthNotServing = 2
)

// registerHealthServer will register a gRPC v1 Health server.
// The serving status will be set to NOT_SERVING.
func (srv *server) registerHealthServer(healthServer *grpc_health_v1.HealthServer) {
	srv.healthz = health.NewServer()
	srv.healthz.SetServingStatus(srv.opts.Name, HealthUnknown)
	if healthServer != nil {
		grpc_health_v1.RegisterHealthServer(srv.grpc, *healthServer)
	} else {
		grpc_health_v1.RegisterHealthServer(srv.grpc, health.NewServer())
	}
	logger.Infof("gRPC health for server '%s' enabled", srv.opts.Name)
}

// setServingStatus of the gRPC server of the health server is enabled.
//
// 0 = UNKNOWN
// 1 = SERVING
// 2 = NOT_SERVING
func (srv *server) setServingStatus(status int32) {
	if srv.opts.HealthServerEnabled {
		srv.healthz.SetServingStatus(srv.opts.Name, grpc_health_v1.HealthCheckResponse_ServingStatus(status))
	}
}
