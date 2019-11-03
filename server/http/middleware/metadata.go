package middleware

import (
	"net/http"

	"github.com/lukasjarosch/genki/metadata"
)

const RequestIDHeaderName = "X-Request-ID"

func Metadata(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		md := metadata.Metadata{}

		ensureRequestId(r, &md)

		ctx = metadata.NewContext(ctx, md)

		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func ensureRequestId(r *http.Request, md *metadata.Metadata) {
	reqId := r.Header.Get(RequestIDHeaderName)
	if reqId != "" {
		(*md)[metadata.RequestIDKey] = reqId
	}
	(*md)[metadata.RequestIDKey] = metadata.NewRequestID()
}