package main

import (
	"log"

	"github.com/lukasjarosch/genki"
	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/broker/amqp"
	"github.com/lukasjarosch/genki/cli"
	"github.com/lukasjarosch/genki/client/http/authz"
	"github.com/lukasjarosch/genki/config"
	events "github.com/lukasjarosch/genki/examples/stringer/internal/broker"
	"github.com/lukasjarosch/genki/examples/stringer/internal/datastore"
	grpc2 "github.com/lukasjarosch/genki/examples/stringer/internal/proto"
	"github.com/lukasjarosch/genki/examples/stringer/internal/stringer"
	"github.com/lukasjarosch/genki/examples/stringer/internal/stringer/middleware"
	"github.com/lukasjarosch/genki/logger"
	"github.com/lukasjarosch/genki/server/grpc"
	"github.com/lukasjarosch/genki/server/http"

	example "github.com/lukasjarosch/genki/examples/stringer/proto"
)

const Service = "stringer"

// init takes care of setting up the CLI flags, parsing and ultimately binding
// them to the configuration. After they are bound, they are globally accessible via the config package.
func init() {
	flags := cli.NewFlagSet(Service)
	flags.Add(logger.Flags, http.Flags, grpc.Flags, amqp.Flags, authz.Flags)
	flags.Parse()
	config.BindFlagSet(flags.Set())
}

func main() {
	if err := logger.NewLogger(config.GetString(logger.LogLevelConfigKey)); err != nil {
		log.Fatal(err.Error())
	}

	amqpBroker := amqp.NewSession(config.GetString(amqp.UrlConfigKey))
	implementation := stringer.NewUseCase(datastore.NewInMem())
	implementation = middleware.NewEventPublisher(events.NewExamplePublisher(amqpBroker))(implementation)
	implementation = middleware.NewExampleMetrics()(implementation)
	implementation = middleware.NewAuthorizer(authz.NewOpenPolicyAgentClient(config.GetString(authz.OpenPolicyAgentUrlConfigKey)))(implementation)
	implementation = middleware.ExampleLogger()(implementation)

	if err := amqpBroker.AddPublisher("test", events.GreetingRenderedTopic); err != nil {
		logger.Warnf("failed to add publisher to exchange '%s' with routing key '%s': %s", "test", events.GreetingRenderedTopic, err)
	}
	if err := amqpBroker.AddSubscription("test", "test-queue", events.GreetingRenderedTopic, events.GreetingHappenedSubscriber(implementation)); err != nil {
		logger.Warnf("failed to add subscription to routing key '%s': %s", err)
	}

	app := genki.NewApplication()
	app.RegisterBroker(amqpBroker)

	// setup gRPC server
	grpcServer := grpc.NewServer(
		grpc.Name(Service),
	)
	example.RegisterExampleServiceServer(grpcServer.Server(), grpc2.NewExampleService(implementation))

	// register servers
	app.AddServer(grpcServer)

	// off we go...
	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

// impl is a quick and dirty handler implementation
type impl struct {
	publisher broker.Publisher
}
