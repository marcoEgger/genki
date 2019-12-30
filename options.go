package genki

import (
	"github.com/spf13/pflag"
)

const (
	DefaultName               = "genki"
	HttpDebugPortConfigKey    = "http-debug-port"
	HttpDebugDisableConfigKey = "http-debug-disable"
)

type Options struct {
	Name                   string
	HttpDebugServerEnabled bool
	HttpDebugServerPort    string
}

func Name(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func DisableDebugHttpServer() Option {
	return func(opts *Options) {
		opts.HttpDebugServerEnabled = false
	}
}

func HttpDebugServerPort(port string) Option {
	return func(opts *Options) {
		if port != "" {
			opts.HttpDebugServerPort = port
		}
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:                   DefaultName,
		HttpDebugServerEnabled: true,
		HttpDebugServerPort:    "3000",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Flags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("genki", pflag.ContinueOnError)
	fs.String(HttpDebugPortConfigKey, "3000", "the port on which the debug http server is listening")
	fs.Bool(HttpDebugDisableConfigKey, false, "disable the http-debug server")
	return fs
}
