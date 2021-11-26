package genki

import (
	"github.com/marcoEgger/genki/broker"
	"github.com/marcoEgger/genki/server"
	"github.com/marcoEgger/genki/server/http"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Application defines the application interface. It's designed to be simple and straightforward.
type Application interface {
	// Name of the application
	Name() string
	// Run the application. This is a blocking call and will only return if the
	// server is shut-down or an error occurred.
	Run(healthServer grpc_health_v1.HealthServer) error
	// Opts returns the current options
	Opts() Options
	// AddServer registers a new server with the application
	AddServer(server server.Server)
	// AddHttpServer registers a new http server with the application
	AddHttpServer(srv http.Server)
	// RegisterBroker registers a message broker implementation (AMQP, NATS, ...)
	RegisterBroker(broker broker.Broker)
}

type Option func(options *Options)

func NewApplication(options ...Option) Application {
	return newService(options...)
}
