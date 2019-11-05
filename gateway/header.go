package gateway

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/lukasjarosch/genki/server/grpc/interceptor"
	"github.com/lukasjarosch/genki/server/http/middleware"
)

// IncomingHeaderMatcher will rewrite HTTP header keys into gRPC header keys.
// All remaining headers are treated with the default policy.
func IncomingHeaderMatcher(key string) (string, bool) {
	switch key {
	case middleware.AccountIDHeaderName:
		return interceptor.AccountIdMetadataKey, true
	case middleware.UserIDHeaderName:
		return interceptor.UserIdMetadataKey, true
	case middleware.RequestIDHeaderName:
		return interceptor.RequestIdMetadataKey, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
