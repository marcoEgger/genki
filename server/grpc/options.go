package grpc

type Options struct {
	Port   string
}

func Port(addr string) Option {
	return func(opts *Options) {
		opts.Port = addr
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Port: "50051",
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
