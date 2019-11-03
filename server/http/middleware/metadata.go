package middleware

import (
	"net/http"

	"github.com/lukasjarosch/genki/metadata"
)

const RequestIDHeaderName = "X-Request-ID"
const AccountIDHeaderName = "X-Account-ID"
const UserIDHeaderName = "X-User-ID"

func Metadata(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		md := metadata.Metadata{}

		ensureRequestID(r, &md)
		findAccountID(r, &md)
		findUserID(r, &md)

		ctx = metadata.NewContext(ctx, md)

		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func ensureRequestID(r *http.Request, md *metadata.Metadata) {
	reqId := r.Header.Get(RequestIDHeaderName)
	if reqId != "" {
		(*md)[metadata.RequestIDKey] = reqId
		return
	}
	(*md)[metadata.RequestIDKey] = metadata.NewRequestID()
}

func findAccountID(r *http.Request, md *metadata.Metadata) {
	accID := r.Header.Get(AccountIDHeaderName)
	if accID != "" {
		(*md)[metadata.AccountIDKey] = accID
	}
}

func findUserID(r *http.Request, md *metadata.Metadata) {
	userID := r.Header.Get(UserIDHeaderName)
	if userID != "" {
		(*md)[metadata.UserIDKey] = userID
	}
}
