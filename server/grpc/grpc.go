package grpc

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"

	"google.golang.org/grpc"
)

type Server interface {
	ListenAndServe(ctx context.Context, wg *sync.WaitGroup, healthServer *grpc_health_v1.HealthServer)
	Server() *grpc.Server
}

const (
	PortConfigKey        = "grpc-port"
	GracePeriodConfigKey = "grpc-grace-period"
)
