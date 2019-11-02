package genki

import (
	"github.com/lukasjarosch/genki/server"
)

const DefaultName = "genki"

type Options struct {
	Name string
	Servers []server.Server
}

func Name(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:   DefaultName,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}