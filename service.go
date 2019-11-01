package genki

import (
	"context"
	"sync"

	"github.com/lukasjarosch/genki/server"
	genki "github.com/lukasjarosch/genki/service"
)

type service struct {
	opts Options
	stopChan <-chan struct{}
	appContext context.Context
	cancel context.CancelFunc
	wg sync.WaitGroup
}

func newService(opts ...Option) *service {
	options := newOptions(opts...)

	svc :=  &service{
		opts:     options,
		stopChan: genki.NewSignalHandler(),
		wg: sync.WaitGroup{},
	}

	svc.appContext, svc.cancel = context.WithCancel(context.Background())

	return svc
}

func (svc *service) Name() string {
	return svc.opts.Name
}

func (svc *service) Run() error {
	defer svc.cancel()

	// start all registered servers in a goroutine
	for _, srv := range svc.opts.Servers {
		svc.wg.Add(1)
		go srv.ListenAndServe(svc.appContext, &svc.wg)
	}

	// wait for signal handler to fire and shutdown
	<-svc.stopChan
	svc.cancel()
	svc.wg.Wait()

	return nil
}

func (svc *service) AddServer(srv server.Server)  {
	svc.opts.Servers = append(svc.opts.Servers, srv)
}

func (svc *service) Opts() Options {
	return svc.opts
}
