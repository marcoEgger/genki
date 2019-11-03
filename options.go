package genki

import (
	"github.com/lukasjarosch/genki/server"
)

const DefaultName = "genki"
const DefaultDebugHttpServerEnabled = true

type Options struct {
	Name string
	Servers []server.Server
	DebugHtpServerEnabled bool
}

func Name(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func DisableDebugHttpServer() Option {
	return func(opts *Options) {
		opts.DebugHtpServerEnabled = false
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:   DefaultName,
		DebugHtpServerEnabled: DefaultDebugHttpServerEnabled,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}