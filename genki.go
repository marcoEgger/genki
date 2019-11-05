package genki

import (
	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/server"
)

// Application defines the application interface. It's designed to be simple and straightforward.
type Application interface {
	// Name of the application
	Name() string
	// Run the application. This is a blocking call and will only return if the
	// server is shut-down or an error occurred.
	Run() error
	// Opts returns the current options
	Opts() Options
	// AddServer registers a new server with the application
	AddServer(server server.Server)
	// RegisterBroker registers a message broker implementation (AMQP, NATS, ...)
	RegisterBroker(broker broker.Broker)
}

type Option func(options *Options)

func NewApplication(options ...Option) Application {
	return newService(options...)
}
