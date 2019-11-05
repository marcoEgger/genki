package stringer

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lukasjarosch/genki/examples/stringer/internal/models"
)

var (
	greetingsRendered = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "example",
		Subsystem:   "domain",
		Name:        "greetings_rendered_count",
		Help:        "the amount of greetings rendered",
		ConstLabels: nil,
	}, []string{"name"})
	domainErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   "example",
		Subsystem:   "domain",
		Name:        "business_error_count",
		Help:        "the amount of errors returned from the business logic",
	})
)

type exampleMetrics struct {
	next Service
}

func NewExampleMetrics() Middleware {
	_ = prometheus.Register(greetingsRendered)
	_ = prometheus.Register(domainErrors)

	return func(next Service) Service {
		return &exampleMetrics{next: next}
	}
}

func (e exampleMetrics) Hello(ctx context.Context, name string) (greeting *models.Greeting, err error) {
	defer func() {
		if err != nil {
		    domainErrors.Inc()
		}
	}()
	return e.next.Hello(ctx, name)
}

func (e exampleMetrics) Render(ctx context.Context, greeting models.Greeting) (err error) {
	defer func() {
		if err != nil {
			domainErrors.Inc()
			return
		}
		greetingsRendered.With(prometheus.Labels{
			"name": greeting.Name,
		}).Inc()
	}()
	return e.next.Render(ctx, greeting)
}
