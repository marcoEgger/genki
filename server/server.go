package server

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"sync"
)

type Server interface {
	ListenAndServe(ctx context.Context, wg *sync.WaitGroup, healthServer *grpc_health_v1.HealthServer)
}

