package http

import (
	"net/http"
	"time"

	"github.com/spf13/pflag"

	"github.com/lukasjarosch/genki/config"
)

const DefaultPort = "8080"
const DefaultGracePeriod = 3 * time.Second

type Options struct {
	Port                string
	ShutdownGracePeriod time.Duration
	Handler http.Handler
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

func Handler(handler http.Handler) Option {
	return func(opts *Options) {
		opts.Handler = handler
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Port:                DefaultPort,
		ShutdownGracePeriod: DefaultGracePeriod,
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
