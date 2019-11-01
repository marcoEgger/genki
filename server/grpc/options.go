package grpc

import "time"

type Options struct {
	Port                string
	ShutdownGracePeriod time.Duration
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

func newOptions(opts ...Option) Options {
	opt := Options{
		Port: "50051",
		ShutdownGracePeriod: 3 * time.Second,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
