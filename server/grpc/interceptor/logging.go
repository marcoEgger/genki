package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/lukasjarosch/genki/logger"
)

func UnaryServerLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := logger.WithMetadata(ctx)
		log.Infof("incoming unary request to '%s'", info.FullMethod)
		defer func(started time.Time) {
			log = logger.WithFields(logger.Fields{
				"took": time.Since(started),
			})
			if err != nil {
				log.Infof("finished unary request to '%s', err=%s", info.FullMethod, err)
				return
			}
			log.Infof("finished unary request to '%s'", info.FullMethod)
		}(time.Now())

		return handler(ctx, req)
	}
}

func UnaryClientLogging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		log := logger.WithMetadata(ctx)
		log.Infof("client call '%s' to server '%s'", method, cc.Target())
		defer func(started time.Time) {
			log = logger.WithFields(logger.Fields{
				"took": time.Since(started),
			})
			if err != nil {
				log.Infof("client call to '%s' (server=%s) failed: %s", method, cc.Target(), err)
			} else {
				log.Infof("client request to '%s' was successfully handled by server '%s'", method, cc.Target())
			}
		}(time.Now())

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
