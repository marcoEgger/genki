package grpc

import (
	"context"
	"sync"

	"github.com/lukasjarosch/genki/logger"
)

type Server struct {
	opts Options
}

func NewServer(opts ...Option) *Server {
	options := newOptions(opts...)

	return &Server{opts: options}
}

func (srv *Server) ListenAndServe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Infof("gRPC server running on port %s", srv.opts.Port)

	<-ctx.Done()

	// TODO: shutdown
}

type Option func (*Options)
