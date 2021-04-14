package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	md "github.com/marcoEgger/genki/metadata"
)

const RequestIdMetadataKey = "requestId"
const AccountIdMetadataKey = "accountId"
const UserIdMetadataKey = "userId"

func UnaryServerMetadata() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		meta := md.Metadata{}
		ensureRequestId(ctx, &meta)
		findAccountId(ctx, &meta)
		findUserId(ctx, &meta)
		ctx = md.NewContext(ctx, meta)

		return handler(ctx, req)
	}
}

func UnaryClientMetadata() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		meta := md.Metadata{}
		ensureRequestId(ctx, &meta)
		findAccountId(ctx, &meta)
		findUserId(ctx, &meta)
		ctx = md.NewContext(ctx, meta)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func ensureRequestId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		requestID := header.Get(RequestIdMetadataKey)
		if len(requestID) > 0 {
			(*meta)[md.RequestIDKey] = requestID[0]
			return
		}
		(*meta)[md.RequestIDKey] = md.NewRequestID()
	}
}

func findAccountId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		accountID := header.Get(AccountIdMetadataKey)
		if len(accountID) > 0 {
			(*meta)[md.AccountIDKey] = accountID[0]
		}
	}
}

func findUserId(ctx context.Context, meta *md.Metadata) {
	if header, ok := metadata.FromIncomingContext(ctx); ok {
		userID := header.Get(UserIdMetadataKey)
		if len(userID) > 0 {
			(*meta)[md.UserIDKey] = userID[0]
		}
	}
}
