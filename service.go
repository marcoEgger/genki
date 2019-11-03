package genki

import (
	"context"
	"fmt"
	"sync"

	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/cli"
	"github.com/lukasjarosch/genki/logger"
	"github.com/lukasjarosch/genki/server"
	"github.com/lukasjarosch/genki/server/http"
	genki "github.com/lukasjarosch/genki/service"
)

type service struct {
	servers    []server.Server
	broker     broker.Broker
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
		flags:    cli.NewFlagSet(options.Name),
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

	if svc.broker != nil {
		if err := svc.broker.Declare(); err != nil {
			return fmt.Errorf("failed to declare broker settings: %s", err)
		}
		svc.wg.Add(1)
		go svc.broker.Consume(&svc.wg)
	}

	// start all registered servers in a goroutine
	for _, srv := range svc.servers {
		svc.wg.Add(1)
		go srv.ListenAndServe(svc.appContext, &svc.wg)
	}

	// wait for signal handler to fire and shutdown
	<-svc.stopChan
	logger.Info("received OS signal: service is shutting down")
	if svc.broker != nil {
		svc.broker.Shutdown()
	}
	svc.cancel()
	svc.wg.Wait()

	return nil
}

// AddServer registers a new server with the service
func (svc *service) AddServer(srv server.Server) {
	svc.servers = append(svc.servers, srv)
}

func (svc *service) RegisterBroker(broker broker.Broker) {
	svc.broker = broker
}

func (svc *service) AddSubscription(exchangeName, queueName, routingKey string, handler broker.Subscriber) error {
	if svc.broker == nil {
		return fmt.Errorf("no broker configured, cannot continue")
	}
	return svc.broker.AddSubscription(exchangeName, queueName, routingKey, handler)
}

func (svc *service) AddPublisher(exchangeName, routingKey string) error {
	if svc.broker == nil {
		return fmt.Errorf("no broker configured, cannot continue")
	}
	return svc.broker.AddPublisher(exchangeName, routingKey)
}

// Opts returns the internal options
func (svc *service) Opts() Options {
	return svc.opts
}
