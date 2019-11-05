package service

import (
	"context"
	"time"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
	"github.com/lukasjarosch/genki/logger"
)

type exampleLogger struct {
	next Service
}

func ExampleLogger() Middleware {
	return func(next Service) Service {
		return &exampleLogger{next: next}
	}
}

func (e exampleLogger) Hello(ctx context.Context, name string) (greeting *models.Greeting, err error) {
	log := logger.WithMetadata(ctx)
	log = log.WithFields(logger.Fields{
		"greeting.name": name,
	})
	log.Infof("call to Hello started")
	defer func(started time.Time) {
		log = log.WithFields(logger.Fields{
			"took": time.Since(started),
		})
		if err != nil {
			log.Infof("call to 'Hello' finished error err=%v", err)
			return
		}
		log = log.WithFields(logger.Fields{
			"greeting.name": greeting.Name,
			"greeting.template": greeting.Template,
		})
		log.Infof("call to 'Hello' finished without error")
	}(time.Now())
	return e.next.Hello(ctx, name)
}

func (e exampleLogger) Render(ctx context.Context, greeting models.Greeting) (err error) {
	log := logger.WithMetadata(ctx)
	log = log.WithFields(logger.Fields{
		"greeting.name": greeting.Name,
		"greeting.template": greeting.Template,
	})
	log.Infof("call to Render started")
	defer func(started time.Time) {
		log = log.WithFields(logger.Fields{
			"took": time.Since(started),
		})
		if err != nil {
			log.Infof("call to render finished with error: %v", err)
			return
		}
		log.Infof("call to render finished without error")
	}(time.Now())
	return e.next.Render(ctx, greeting)
}
