package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/lukasjarosch/genki/logger"
)

type Server struct {
	opts Options
	grpc *grpc.Server
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	return &Server{opts: options}
}

func (srv *Server) ListenAndServe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	srv.grpc = grpc.NewServer()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", srv.opts.Port))
	if err != nil {
		logger.Fatalf("failed to listen on tcp port '%s': %s", srv.opts.Port, err.Error())
	}

	go func() {
		logger.Infof("gRPC server running on port %s", srv.opts.Port)
		if err := srv.grpc.Serve(listener); err != nil {
			logger.Errorf("gRPC server crashed: %s", err.Error())
		}
	}()

	<-ctx.Done()
	srv.shutdown()
}

func (srv *Server) shutdown() {
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
