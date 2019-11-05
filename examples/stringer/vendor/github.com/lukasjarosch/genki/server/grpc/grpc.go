package grpc

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

type Server interface {
	ListenAndServe(ctx context.Context, wg *sync.WaitGroup)
	Server() *grpc.Server
}

const (
	PortConfigKey        = "grpc-port"
	GracePeriodConfigKey = "grpc-grace-period"
)
