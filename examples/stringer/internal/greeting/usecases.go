package greeting

import (
	"context"

	"github.com/lukasjarosch/genki/broker"
	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
	example "github.com/lukasjarosch/genki/examples/stringer/proto"
	"github.com/lukasjarosch/genki/logger"
)

type ServiceImplementation struct {
	publisher broker.Publisher
}

func NewServiceImplementation(publisher broker.Publisher) Service {
	return &ServiceImplementation{publisher: publisher}
}

func (svc *ServiceImplementation) Hello(ctx context.Context, name string) (greeting *models.Greeting, err error) {
	log := logger.WithMetadata(ctx)

	if err := svc.publisher.Publish("some.key", &example.Greeting{}); err != nil {
		log.Warnf("failed to publish to '%s': %s", "some.key", err)
	}
	log.Infof("rendered greeting: %s", greeting)

	greeting = &models.Greeting{
		Name: name,
	}
	if err := greeting.Validate(); err != nil {
		return nil, models.ErrGreetingEmptyName
	}
	greeting.Template = "Ohai there, %s"
	greeting.Rendered = greeting.Render()

	return greeting, nil
}
