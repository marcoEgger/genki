package middleware

import (
	"context"

	"github.com/marcoEgger/genki/client/http/authz"
	"github.com/marcoEgger/genki/examples/stringer/internal/models"
	"github.com/marcoEgger/genki/examples/stringer/internal/stringer"
	"github.com/marcoEgger/genki/logger"
)

type authorizer struct {
	next stringer.Service
	auth authz.Authorizer
}

func NewAuthorizer(auth authz.Authorizer) stringer.Middleware {
	return func(next stringer.Service) stringer.Service {
		return &authorizer{next: next, auth: auth}
	}
}

// Hello will publish an event if the use-case executed without an error
func (a authorizer) Hello(ctx context.Context, name string) (greeting *models.Greeting, err error) {
	log := logger.WithMetadata(ctx)

	err = a.auth.Authorize(ctx, "no", "get", nil)
	if err != nil {
		log.Infof("request to %s is not authorized, reject", "Hello")
		return nil, models.Unauthorized
	}
	log.Infof("request to %s is authorized, continue", "Hello")

	return a.next.Hello(ctx, name)
}

// Render: nothing to do
func (a authorizer) Render(ctx context.Context, greeting models.Greeting) (err error) {
	log := logger.WithMetadata(ctx)

	err = a.auth.Authorize(ctx, "no", "get", nil)
	if err != nil {
		log.Infof("request to %s is not authorized, reject", "Render")
		return  models.Unauthorized
	}
	log.Infof("request to %s is authorized, continue", "Render")

	return a.next.Render(ctx, greeting)
}
