package http

import (
	"time"

	"github.com/spf13/pflag"

	"github.com/lukasjarosch/genki/config"
)

const DefaultPort = "8080"
const DefaultGracePeriod = 3 * time.Second
const DefaultName = "default"
const DefaultLoggingMiddlewareEnabled = true
const DefaultRequestIdMiddlewareEnabled = true

type Options struct {
	Port                       string
	Name                       string
	ShutdownGracePeriod        time.Duration
	LoggingMiddlewareEnabled   bool
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

func Name(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func DisableLoggingMiddleware() Option {
	return func(opts *Options) {
		opts.LoggingMiddlewareEnabled = false
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Port:                       DefaultPort,
		ShutdownGracePeriod:        DefaultGracePeriod,
		Name:                       DefaultName,
		LoggingMiddlewareEnabled:   DefaultLoggingMiddlewareEnabled,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Flags is a convenience function to quickly add the HTTP server options as CLI flags
// Implements the cli.FlagProvider type
func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("http-server", pflag.ContinueOnError)

	fs.String(
		config.HttpPort,
		DefaultPort,
		"the port on which the HTTP server is listening on",
	)
	fs.Duration(
		config.HttpGracePeriod,
		DefaultGracePeriod,
		"grace period after which the server shutdown is terminated",
	)

	return fs
}
