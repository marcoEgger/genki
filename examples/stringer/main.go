package main

import (
	"log"

	amqp2 "github.com/streadway/amqp"

	"github.com/lukasjarosch/genki"
	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/broker/amqp"
	"github.com/lukasjarosch/genki/cli"
	"github.com/lukasjarosch/genki/config"
	"github.com/lukasjarosch/genki/examples/stringer/internal/greeting"
	grpc2 "github.com/lukasjarosch/genki/examples/stringer/internal/grpc"
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
	flags.Add(logger.Flags, http.Flags, grpc.Flags, amqp.Flags)
	flags.Parse()
	config.BindFlagSet(flags.Set())
}

func main() {
	if err := logger.NewLogger(config.GetString(logger.LogLevelConfigKey)); err != nil {
		log.Fatal(err.Error())
	}

	// setup broker
	amqpBroker := amqp.NewSession(config.GetString(amqp.UrlConfigKey))
	h1 := func(delivery interface{}) {
		event := delivery.(amqp2.Delivery)
		logger.Info("OHAI")
		_ = event.Ack(false)
	}
	if err := amqpBroker.AddPublisher("test", "some.key"); err != nil {
		logger.Warnf("failed to add publisher to exchange '%s' with routing key '%s': %s", "test", "some.key", err)
	}
	if err := amqpBroker.AddSubscription("test", "test-queue", "some.key", h1); err != nil {
		logger.Warnf("failed to add subscription to routing key '%s': %s", err)
	}

	// service layer
	implementation := greeting.NewServiceImplementation(amqpBroker)

	// transport layer
	app := genki.NewService()
	app.RegisterBroker(amqpBroker)

	// setup gRPC server
	grpcServer := grpc.NewServer(
		grpc.Name(Service),
		grpc.Port(config.GetString(grpc.PortConfigKey)),
		grpc.ShutdownGracePeriod(config.GetDuration(grpc.GracePeriodConfigKey)),
		grpc.EnableHealthServer(Service),
	)
	example.RegisterExampleServiceServer(grpcServer.Server(), grpc2.NewExampleService(implementation))

	// register servers
	app.AddServer(grpcServer)
	app.AddServer(http.NewServer())
	app.AddServer(http.NewServer(http.Port("12345"), http.Name("pizza")))

	// off we go...
	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

// impl is a quick and dirty handler implementation
type impl struct {
	publisher broker.Publisher
}
