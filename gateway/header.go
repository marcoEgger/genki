package gateway

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/lukasjarosch/genki/server/grpc/interceptor"
	"github.com/lukasjarosch/genki/server/http/middleware"
)

// IncomingHeaderMatcher will rewrite HTTP header keys into gRPC header keys.
// All remaining headers are treated with the default policy.
// The header comparison is case-insensitive.
func IncomingHeaderMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case strings.ToLower(middleware.AccountIDHeaderName):
		return interceptor.AccountIdMetadataKey, true
	case strings.ToLower(middleware.UserIDHeaderName):
		return interceptor.UserIdMetadataKey, true
	case strings.ToLower(middleware.RequestIDHeaderName):
		return interceptor.RequestIdMetadataKey, true
	case strings.ToLower(middleware.RequestIDGatewayHeaderName):
		return interceptor.RequestIdMetadataKey, true
	case strings.ToLower(middleware.EmailHeaderName):
		return interceptor.EmailMetadataKey, true
	case strings.ToLower(middleware.FirstNameHeaderName):
		return interceptor.FirstNameMetadataKey, true
	case strings.ToLower(middleware.LastNameHeaderName):
		return interceptor.LastNameMetadataKey, true
	case strings.ToLower(middleware.TypeHeaderName):
		return interceptor.TypeMetadataKey, true
	case strings.ToLower(middleware.SubTypeHeaderName):
		return interceptor.SubTypeMetadataKey, true
	case strings.ToLower(middleware.RolesHeaderName):
		return interceptor.RolesMetadataKey, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
