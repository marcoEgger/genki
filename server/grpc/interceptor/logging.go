package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/lukasjarosch/genki/logger"
)

func Logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := logger.WithContext(ctx)
		log.Infof("incoming unary request to '%s'", info.FullMethod)
		defer func(started time.Time) {
			log.Infof("finished unary request to '%s' (took %v)", info.FullMethod, time.Since(started))
		}(time.Now())

		return handler(ctx, req)
	}
}
