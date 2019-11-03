package main

import (
	"context"
	"log"

	"github.com/spf13/pflag"
	amqp2 "github.com/streadway/amqp"

	"github.com/lukasjarosch/genki"
	"github.com/lukasjarosch/genki/broker/amqp"
	"github.com/lukasjarosch/genki/cli"
	"github.com/lukasjarosch/genki/config"
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
	flags.Add(logger.Flags)
	flags.Add(http.Flags)
	flags.Add(grpc.Flags)
	flags.Add(amqp.Flags)
	flags.Add(Flags)
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


	app := genki.NewService()
	app.RegisterBroker(amqpBroker)

	// setup gRPC server
	grpcServer := grpc.NewServer(
		grpc.Name(Service),
		grpc.Port(config.GetString(grpc.PortConfigKey)),
		grpc.ShutdownGracePeriod(config.GetDuration(grpc.GracePeriodConfigKey)),
		grpc.EnableHealthServer(Service),
	)
	example.RegisterExampleServiceServer(grpcServer.Server(), &impl{
		publisher:amqpBroker,
	})

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
	publisher amqp.Publisher
}

func (i *impl) Hello(ctx context.Context, request *example.HelloRequest) (*example.HelloResponse, error) {

	if err := i.publisher.Publish("some.key", &example.Greeting{}); err != nil {
		logger.Warnf("failed to publish to '%s': %s", "some.key", err)
	}

	return &example.HelloResponse{
		Greeting:             &example.Greeting{
			Template:             "Ohai there %s",
			Name:                 "Lukas",
			Rendered:             "Ohai there Lukas",
		},
	}, nil
}

// this could be somewhere in your business logic. That's an easy way to configure it.
func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("http-server", pflag.ContinueOnError)

	fs.String(
		"my-own-flag",
		"foobar",
		"something something description",
	)

	return fs
}