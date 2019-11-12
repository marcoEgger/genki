package middleware

import (
	"context"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
	"github.com/lukasjarosch/genki/examples/stringer/internal/stringer"
	"github.com/lukasjarosch/genki/logger"
)

type eventPublisher struct {
	next stringer.Service
	pub  stringer.Publisher
}

func NewEventPublisher(publisher stringer.Publisher) stringer.Middleware {
	return func(next stringer.Service) stringer.Service {
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
