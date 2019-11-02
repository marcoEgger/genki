package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	metadata2 "github.com/lukasjarosch/genki/server/grpc/metadata"
)

const RequestIdKey = metadata2.RequestID

func RequestId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if meta, ok := metadata.FromIncomingContext(ctx); ok {

			requestID := meta.Get(RequestIdKey)
			if len(requestID) > 0 {
				ctx = context.WithValue(ctx, RequestIdKey, requestID)
				return handler(ctx, req)
			}

			newRequestID := newRequestID()
			meta.Append(RequestIdKey, newRequestID)
			ctx = metadata.NewIncomingContext(ctx, meta)
			ctx = context.WithValue(ctx, RequestIdKey, newRequestID)
			return handler(ctx, req)
		}

		newRequestID := newRequestID()
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(RequestIdKey, newRequestID))
		ctx = context.WithValue(ctx, RequestIdKey, newRequestID)
		return handler(ctx, req)
	}

}

func newRequestID() string {
	return uuid.New().String()
}
