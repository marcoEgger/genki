package service

import (
	"context"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
	"github.com/lukasjarosch/genki/logger"
)

type eventPublisher struct {
	next Service
	pub  Publisher
}

func NewEventPublisher(publisher Publisher) Middleware {
	return func(next Service) Service {
		return &eventPublisher{next: next, pub: publisher}
	}
}

// Hello will publish an event if the use-case executed without an error
func (e eventPublisher) Hello(ctx context.Context, name string) (greeting *models.Greeting, err error) {
	log := logger.WithMetadata(ctx)
	defer func() {
		if err == nil {
			if err := e.pub.GreetingRendered(greeting); err != nil {
				log.Warnf("failed to publish Greeting with routing key '%s': %s", "some.key", err)
			}
		}
	}()
	return e.next.Hello(ctx, name)
}

// Render: nothing to do
func (e eventPublisher) Render(ctx context.Context, greeting models.Greeting) (err error) {
	return e.next.Render(ctx, greeting)
}
