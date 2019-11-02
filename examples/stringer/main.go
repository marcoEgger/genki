package main

import (
	"context"
	"log"

	"github.com/spf13/pflag"

	"github.com/lukasjarosch/genki"
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
	flags.Add(Flags)
	flags.Parse()
	config.BindFlagSet(flags.Set())
}

func main() {
	if err := logger.NewLogger(config.GetString(config.LogLevel)); err != nil {
		log.Fatal(err.Error())
	}

	app := genki.NewService()

	// setup gRPC server
	grpcServer := grpc.NewServer(
		grpc.Port(config.GetString(config.GrpcPort)),
		grpc.ShutdownGracePeriod(config.GetDuration(config.GrpcGracePeriod)),
		grpc.EnableHealthServer(Service),
	)
	example.RegisterExampleServiceServer(grpcServer.Server(), &impl{})

	// setup HTTP server
	httpServer := http.NewServer(
		http.Port("3000"),
		http.Handler(nil),
	)

	// register servers
	app.AddServer(grpcServer)
	app.AddServer(httpServer)

	// run application
	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

// impl is a quick and dirty handler implementation
type impl struct {
}

func (i *impl) Hello(ctx context.Context, request *example.HelloRequest) (*example.HelloResponse, error) {
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