package gateway

import (
	"context"
	"github.com/golang/protobuf/proto"
	"net/http"

)

type Options struct {
	ResponseInterceptor ResponseInterceptorFunc
}
type ResponseInterceptorFunc func(context.Context, http.ResponseWriter, proto.Message) error

func ResponseInterceptor(interceptor ResponseInterceptorFunc) Option {
	return func(opt *Options) {
		opt.ResponseInterceptor = interceptor
	}
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		ResponseInterceptor: nil,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

type Option func(*Options)
