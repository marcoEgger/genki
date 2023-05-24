package gateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type Options struct {
	ResponseInterceptor ResponseInterceptorFunc
	ServeMuxOpts        []runtime.ServeMuxOption
}
type ResponseInterceptorFunc func(context.Context, http.ResponseWriter, proto.Message) error

func ResponseInterceptor(interceptor ResponseInterceptorFunc) Option {
	return func(opt *Options) {
		opt.ResponseInterceptor = interceptor
	}
}

func WithServeMuxOptions(opts ...runtime.ServeMuxOption) Option {
	return func(opt *Options) {
		opt.ServeMuxOpts = append(opt.ServeMuxOpts, opts...)
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
