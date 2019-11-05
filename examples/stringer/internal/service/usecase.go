package service

import (
	"context"
	"strings"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
	"github.com/lukasjarosch/genki/logger"
)

type useCase struct {
	exampleRepo models.GreetingRepository
}

func NewUseCase(repo models.GreetingRepository) Service {
	return &useCase{exampleRepo: repo}
}

func (svc *useCase) Hello(ctx context.Context, name string) (greeting *models.Greeting, err error) {
	log := logger.WithMetadata(ctx)

	greeting, _ = svc.exampleRepo.FindGreetingByName(name)
	if err := greeting.Validate(); err != nil {
		return nil, models.ErrGreetingEmptyName
	}
	greeting.Rendered = greeting.Render()
	log.Infof("rendered greeting: %s", greeting)

	return greeting, nil
}

func (svc *useCase) Render(ctx context.Context, greeting models.Greeting) error {
	logger.Info(strings.Repeat("YO", len(greeting.Render())/2))
	logger.Infof("%s", greeting.Render())
	logger.Info(strings.Repeat("YO", len(greeting.Render())/2))
	return nil
}
