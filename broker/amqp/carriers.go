package amqp

import (
	"context"
	"go.opentelemetry.io/otel"
)

type HeadersCarrier map[string]interface{}

func (a HeadersCarrier) Get(key string) string {
	v, ok := a[key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (a HeadersCarrier) Set(key string, value string) {
	a[key] = value
}

func (a HeadersCarrier) Keys() []string {
	i := 0
	r := make([]string, len(a))

	for k, _ := range a {
		r[i] = k
		i++
	}

	return r
}

// InjectAMQPHeaders injects the tracing from the context into the header map
func InjectAMQPHeaders(ctx context.Context) map[string]interface{} {
	h := make(HeadersCarrier)
	otel.GetTextMapPropagator().Inject(ctx, h)
	return h
}

// ExtractAMQPHeaders extracts the tracing from the header and puts it into the context
func ExtractAMQPHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, HeadersCarrier(headers))
}
