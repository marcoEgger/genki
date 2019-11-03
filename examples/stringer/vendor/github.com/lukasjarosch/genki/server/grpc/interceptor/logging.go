package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/lukasjarosch/genki/logger"
)

func Logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := logger.WithContext(ctx)
		log.Infof("incoming unary request to '%s': %v", info.FullMethod, ctx)
		defer func() {

		}()

		return handler(ctx, req)
	}
}