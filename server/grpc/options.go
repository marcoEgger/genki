package grpc

import (
	"time"

	"github.com/spf13/pflag"

	"github.com/lukasjarosch/genki/config"
)

const DefaultPort = "50051"
const DefaultGracePeriod = 3 * time.Second
const DefaultHealthEnabled = false
const DefaultLoggingInterceptor = true
const DefaultRequestIdInterceptor = true

type enabledInterceptor struct {
	logging   bool
	requestId bool
}

type Options struct {
	Port                    string
	ShutdownGracePeriod     time.Duration
	HealthServerEnabled     bool
	serviceName             string // only set if EnableHealthServer is called
	enabledUnaryInterceptor enabledInterceptor
}

func Port(addr string) Option {
	return func(opts *Options) {
		opts.Port = addr
	}
}

func ShutdownGracePeriod(duration time.Duration) Option {
	return func(opts *Options) {
		opts.ShutdownGracePeriod = duration
	}
}

func EnableHealthServer(serviceName string) Option {
	return func(opts *Options) {
		opts.HealthServerEnabled = true
		opts.serviceName = serviceName
	}
}

func DisableLoggingInterceptor() Option {
	return func(opts *Options) {
		opts.enabledUnaryInterceptor.logging = false
	}
}

func DisableRequestIdInterceptor() Option {
	return func(opts *Options) {
		opts.enabledUnaryInterceptor.requestId = false
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Port:                DefaultPort,
		ShutdownGracePeriod: DefaultGracePeriod,
		HealthServerEnabled: DefaultHealthEnabled,
		enabledUnaryInterceptor: enabledInterceptor{
			logging:   DefaultLoggingInterceptor,
			requestId: DefaultRequestIdInterceptor,
		},
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Flags is a convenience function to quickly add the gRPC server options as CLI flags
// Implements the cli.FlagProvider type
func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("grpc-server", pflag.ContinueOnError)

	fs.String(
		config.GrpcPort,
		DefaultPort,
		"the port on which the gRPC server is listening on",
	)
	fs.Duration(
		config.GrpcGracePeriod,
		DefaultGracePeriod,
		"grace period after which the server shutdown is terminated",
	)

	return fs
}
