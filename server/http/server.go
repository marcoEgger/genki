package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/lukasjarosch/genki/logger"
	"github.com/lukasjarosch/genki/server/http/middleware"
)

type server struct {
	opts       Options
	mux        *http.ServeMux
	httpServer *http.Server
}

func NewServer(opts ...Option) Server {
	options := newOptions(opts...)

	return &server{
		opts: options,
		mux:  http.NewServeMux(),
	}
}

func (srv *server) Handle(endpoint string, handler http.Handler) {
	if srv.opts.LoggingMiddlewareEnabled {
		handler = middleware.Logging(handler)
	}

	srv.mux.Handle(endpoint, handler)
}

func (srv *server) ListenAndServe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", srv.opts.Port),
		Handler: srv.mux,
	}

	// serve
	go func() {
		logger.Infof("HTTP server '%s' running on port %s", srv.opts.Name, srv.opts.Port)
		if err := srv.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP server crashed: %s", err.Error())
			return
		}
	}()

	<-ctx.Done()
}

func (srv *server) shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), srv.opts.ShutdownGracePeriod)
	defer cancel()
	if err := srv.httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Warnf("gRPC graceful shutdown timed-out: %s", err.Error())
	} else {
		logger.Info("HTTP server shut-down gracefully")
	}
}

type Option func(*Options)
