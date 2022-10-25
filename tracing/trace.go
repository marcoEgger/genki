package tracing

import (
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.12.0"
)

//goland:noinspection GoUnusedExportedFunction
func InitTracing(service string, namespace string, url string) error {
	// Create jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	traceProvider := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
		semconv.ServiceNamespaceKey.String(namespace),
	)))
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(b3.New())
	return nil
}
