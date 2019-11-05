package greeting

import (
	"context"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
)

type Service interface {
	Hello(ctx context.Context, name string) (greeting *models.Greeting, err error)
}

type Middleware func(svc Service) Service

type Repository interface {
	GetHelloTemplate(name string) (format string, err error)
}
