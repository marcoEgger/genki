package genki

import (
	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/server"
)

type Service interface {
	// Name of the service
	Name() string
	// Run the service. This is a blocking call and will only return if the
	// server is shut-down or an error occurred.
	Run() error
	// Opts returns the current options
	Opts() Options
	// AddServer registers a new server with the application
	AddServer(server server.Server)
	// RegisterBroker registers a message broker implementation (AMQP, NATS, ...)
	RegisterBroker(broker broker.Broker)

	broker.SubscriptionCreator
	broker.PublishCreator
}

type Option func(options *Options)

func NewService(options ...Option) Service {
	return newService(options...)
}