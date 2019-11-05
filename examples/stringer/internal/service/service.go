package service

import (
	"context"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
)

type Service interface {
	Hello(ctx context.Context, name string) (greeting *models.Greeting, err error)
	Render(ctx context.Context, greeting models.Greeting) error
}

type Middleware func(svc Service) Service

type Repository interface {
	GetHelloTemplate(name string) (format string, err error)
}

type Publisher interface {
	GreetingRendered(greeting *models.Greeting) error
}
