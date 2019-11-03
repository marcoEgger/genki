package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/lukasjarosch/genki/logger"
	"github.com/lukasjarosch/genki/server/grpc/interceptor"
)

type server struct {
	opts Options
	grpc *grpc.Server

	// only set if the gRPC health server is enabled
	healthz *health.Server
}

func NewServer(opts ...Option) Server {
	options := newOptions(opts...)
	srv := &server{opts: options}

	var unaryInterceptors []grpc.UnaryServerInterceptor
	if srv.opts.enabledUnaryInterceptor.requestId {
		logger.Debugf("grpc interceptor enabled: RequestId")
		unaryInterceptors = append(unaryInterceptors, interceptor.RequestId())
	}
	if srv.opts.enabledUnaryInterceptor.logging {
		logger.Debugf("grpc interceptor enabled: Logging")
		unaryInterceptors = append(unaryInterceptors, interceptor.Logging())
	}

	srv.grpc = grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		unaryInterceptors...
	)))

	return  srv
}

// ListenAndServe ties everything together and runs the gRPC server in a separate goroutine.
// The method then blocks until the passed context is cancelled, so this method should also be started
// as goroutine if more work is needed after starting the gRPC server.
func (srv *server) ListenAndServe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	if srv.opts.HealthServerEnabled {
		srv.registerHealthServer()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", srv.opts.Port))
	if err != nil {
		srv.setServingStatus(HealthNotServing)
		logger.Fatalf("gRPC server could not be started: %s", err.Error())
	}



	go func() {
		logger.Infof("gRPC server running on port %s", srv.opts.Port)
		srv.setServingStatus(HealthServing)
		if err := srv.grpc.Serve(listener); err != nil {
			srv.setServingStatus(HealthNotServing)
			logger.Errorf("gRPC server crashed: %s", err.Error())
			return
		}
	}()

	<-ctx.Done()
	srv.setServingStatus(HealthNotServing)

	srv.shutdown()
}

// Server returns the raw grpc Server instance. It is required to register services.
func (srv *server) Server() *grpc.Server {
	return srv.grpc
}

// shutdown is responsible of gracefully shutting down the gRPC server
// First, GracefulStop() is executed. If the call doesn't return
// until the ShutdownGracePeriod is exceeded, the server is terminated directly.
func (srv *server) shutdown() {
	stopped := make(chan struct{})
	go func() {
		srv.grpc.GracefulStop()
		close(stopped)
	}()
	t := time.NewTicker(srv.opts.ShutdownGracePeriod)
	select {
	case <-t.C:
		logger.Warnf("gRPC graceful shutdown timed-out")
	case <-stopped:
		logger.Info("gRPC server shut-down gracefully")
		t.Stop()
	}
}

type Option func(*Options)
