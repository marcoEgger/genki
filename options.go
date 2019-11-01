package genki

import (
	"github.com/lukasjarosch/genki/server"
)

type Options struct {
	Name string
	Servers []server.Server
}

func Name(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func AddServer(srv server.Server) Option {
	return func(opts *Options) {
		opts.Servers = append(opts.Servers, srv)
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:   "genki",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}