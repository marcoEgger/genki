package gateway

import (
	"context"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"github.com/lukasjarosch/genki/server/grpc/interceptor"
)

type gateway struct {
	ctx         context.Context
	mux         *runtime.ServeMux
	dialOptions []grpc.DialOption
	opts        *Options
}

type Gateway interface {
	HttpMux() *runtime.ServeMux
	GrpcDialOpts() []grpc.DialOption
	Context() context.Context
}

func NewGateway(ctx context.Context, options ...Option) Gateway {
	opts := newOptions(options...)

	// configure gateway runtime options
	serveMuxOpts := []runtime.ServeMuxOption{
		runtime.WithIncomingHeaderMatcher(IncomingHeaderMatcher),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
	}

	// allow external response middleware
	if opts.ResponseInterceptor != nil {
		serveMuxOpts = append(serveMuxOpts, runtime.WithForwardResponseOption(opts.ResponseInterceptor))
	}

	mux := runtime.NewServeMux(serveMuxOpts...)
	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(
			interceptor.UnaryClientLogging(),
			interceptor.UnaryClientPrometheus(),
			interceptor.UnaryClientMetadata(),
		)),
	}

	gw := &gateway{
		ctx:         ctx,
		mux:         mux,
		dialOptions: dialOpts,
		opts:        opts,
	}

	return gw
}

func (gw *gateway) HttpMux() *runtime.ServeMux {
	return gw.mux
}

func (gw *gateway) GrpcDialOpts() []grpc.DialOption {
	return gw.dialOptions
}

func (gw *gateway) Context() context.Context {
	return gw.ctx
}
