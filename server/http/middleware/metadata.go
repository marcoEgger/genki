package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/lukasjarosch/genki/metadata"
)

const RequestHeaderName = "X-Request-ID"

func RequestId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get(RequestHeaderName)
		if reqId != "" {
			ctx := metadata.NewContext(r.Context(), metadata.Metadata{
				metadata.RequestIDKey: reqId,
			})
			r = r.WithContext(ctx)
			handler.ServeHTTP(w, r)
			return
		}

		// no X-Request-ID header, create new requestId
		ctx := metadata.NewContext(r.Context(), metadata.Metadata{
			metadata.RequestIDKey: newRequestID(),
		})
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	return uuid.New().String()
}
