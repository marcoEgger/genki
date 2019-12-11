package middleware

import (
	"encoding/base64"
	"net/http"

	"github.com/lukasjarosch/genki/metadata"
)

const RequestIDHeaderName = "X-Request-ID"
const AccountIDHeaderName = "X-Account-ID"
const UserIDHeaderName = "X-User-ID"
const EmailHeaderName = "X-E-Mail"
const FirstNameHeaderName = "X-First-Name"
const LastNameHeaderName = "X-Last-Name"
const TypeHeaderName = "X-Type"
const SubTypeHeaderName = "X-Sub-Type"
const RolesHeaderName = "X-Roles"
const RequestIDGatewayHeaderName = "eg-request-id"

func Metadata(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		md := metadata.Metadata{}

		ensureRequestID(r, &md)
		findAccountID(r, &md)
		findUserID(r, &md)
		findEmail(r, &md)
		findFirstName(r, &md)
		findLastName(r, &md)
		findType(r, &md)
		findSubType(r, &md)
		findRoles(r, &md)

		ctx = metadata.NewContext(ctx, md)

		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func ensureRequestID(r *http.Request, md *metadata.Metadata) {
	// look for a default requestId header
	reqId := r.Header.Get(RequestIDHeaderName)
	if reqId != "" {
		(*md)[metadata.RequestIDKey] = reqId
		return
	}

	// maybe the gateway provided one
	reqId = r.Header.Get(RequestIDGatewayHeaderName)
	if reqId != "" {
		(*md)[metadata.RequestIDKey] = reqId
		return
	}

	reqID := metadata.NewRequestID()
	(*md)[metadata.RequestIDKey] = reqID
	r.Header.Set(RequestIDHeaderName, reqID)
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

func findEmail(r *http.Request, md *metadata.Metadata) {
	email := r.Header.Get(EmailHeaderName)
	if email != "" {
		(*md)[metadata.EmailKey] = base64.StdEncoding.EncodeToString([]byte(email))
	}
}

func findFirstName(r *http.Request, md *metadata.Metadata) {
	firstName := r.Header.Get(FirstNameHeaderName)
	if firstName != "" {
		(*md)[metadata.FirstNameKey] = base64.StdEncoding.EncodeToString([]byte(firstName))
	}
}

func findLastName(r *http.Request, md *metadata.Metadata) {
	lastName := r.Header.Get(LastNameHeaderName)
	if lastName != "" {
		(*md)[metadata.LastNameKey] = base64.StdEncoding.EncodeToString([]byte(lastName))
	}
}

func findType(r *http.Request, md *metadata.Metadata) {
	typ := r.Header.Get(TypeHeaderName)
	if typ != "" {
		(*md)[metadata.TypeKey] = typ
	}
}

func findSubType(r *http.Request, md *metadata.Metadata) {
	subType := r.Header.Get(SubTypeHeaderName)
	if subType != "" {
		(*md)[metadata.SubTypeKey] = subType
	}
}

func findRoles(r *http.Request, md *metadata.Metadata) {
	subType := r.Header.Get(RolesHeaderName)
	if subType != "" {
		(*md)[metadata.RolesKey] = subType
	}
}
