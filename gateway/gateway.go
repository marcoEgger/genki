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
}

type Gateway interface {
	HttpMux() *runtime.ServeMux
	GrpcDialOpts() []grpc.DialOption
	Context() context.Context
}

func NewGateway(ctx context.Context) Gateway {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(IncomingHeaderMatcher),
	)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(
			interceptor.UnaryClientMetadata(),
			interceptor.UnaryClientPrometheus(),
			interceptor.UnaryClientLogging(),
		)),
	}

	gw := &gateway{
		ctx:         ctx,
		mux:         mux,
		dialOptions: opts,
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

