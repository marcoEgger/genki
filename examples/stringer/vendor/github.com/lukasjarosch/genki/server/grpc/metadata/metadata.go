package metadata

import (
	"context"

	"google.golang.org/grpc/metadata"
)

const (
	RequestID string = "requestId"
	AccountID string = "accountId"
	UserID    string = "userId"
	TraceID   string = "zipkinTraceId"
)

// GetMetadata is a convenience function which can be used in order to not have to import two metadata
// libraries (grpc/metadata and go-godin/metadata)
func GetMetadata(ctx context.Context) (metadata.MD, bool) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md, true
	}
	return nil, false
}

// GetRequestID tries to extract the requestId key from the given context.
func GetRequestID(ctx context.Context) string {

	// try and find it in grpc metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		requestID := md.Get(string(RequestID))
		if len(requestID) > 0 {
			return requestID[0]
		}
	}

	// requestId might also be in the context already (e.g. from an AMQP subscriber which does not have metadata)
	requestId := ctx.Value(string(RequestID))
	if requestId.(string) != "" {
		return requestId.(string)
	}

	return ""
}

// GetAccountID tries to extract the accountId key from the given context.
// If no AccountID exists, an empty string is returned.
func GetAccountID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		accountID := md.Get(string(AccountID))
		if len(accountID) > 0 {
			return accountID[0]
		}
	}
	return ""
}

// GetUserID tries to extract the userId key from the given context.
// If no UserID exists, an empty string is returned
func GetUserID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		userID := md.Get(string(UserID))
		if len(userID) > 0 {
			return userID[0]
		}
	}
	return ""
}

// Has checks whether the passed key exists in the context metadata
func Has(ctx context.Context, key string) bool {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get(string(key))
		if len(vals) > 0 {
			return true
		}
	}
	return false
}

func HasAccountID(ctx context.Context) bool {
	return Has(ctx, AccountID)
}

func HasUserID(ctx context.Context) bool {
	return Has(ctx, UserID)
}
