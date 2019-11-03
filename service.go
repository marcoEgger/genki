package genki

import (
	"context"
	"sync"

	"github.com/lukasjarosch/genki/cli"
	"github.com/lukasjarosch/genki/logger"
	"github.com/lukasjarosch/genki/server"
	"github.com/lukasjarosch/genki/server/http"
	genki "github.com/lukasjarosch/genki/service"
)

type service struct {
	opts       Options
	stopChan   <-chan struct{}
	appContext context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	flags      *cli.FlagSet
}

func newService(opts ...Option) *service {
	options := newOptions(opts...)

	svc := &service{
		opts:     options,
		stopChan: genki.NewSignalHandler(),
		wg:       sync.WaitGroup{},
		flags: cli.NewFlagSet(options.Name),
	}

	svc.appContext, svc.cancel = context.WithCancel(context.Background())

	return svc
}

func (svc *service) Name() string {
	return svc.opts.Name
}

func (svc *service) Run() error {
	defer svc.cancel()

	// add the debug HTTP server if enabled
	if svc.opts.DebugHtpServerEnabled {
		svc.AddServer(http.NewDebugServer())
	}

	// start all registered servers in a goroutine
	for _, srv := range svc.opts.Servers {
		svc.wg.Add(1)
		go srv.ListenAndServe(svc.appContext, &svc.wg)
	}

	// wait for signal handler to fire and shutdown
	<-svc.stopChan
	logger.Info("received OS signal: service is shutting down")
	svc.cancel()
	svc.wg.Wait()

	return nil
}


// AddServer registers a new server with the service
func (svc *service) AddServer(srv server.Server) {
	svc.opts.Servers = append(svc.opts.Servers, srv)
}

// Opts returns the internal options
func (svc *service) Opts() Options {
	return svc.opts
}
