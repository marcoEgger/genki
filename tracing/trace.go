package tracing

import (
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.12.0"
	"strings"
)

type CustomSampler struct{}

func (cs *CustomSampler) ShouldSample(p trace.SamplingParameters) trace.SamplingResult {
	if strings.Contains(p.Name, "ping") {
		return trace.SamplingResult{Decision: trace.Drop}
	}
	return trace.SamplingResult{Decision: trace.RecordAndSample}
}

func (cs *CustomSampler) Description() string {
	return "CustomSampler"
}

//goland:noinspection GoUnusedExportedFunction
func InitTracing(service string, namespace string, url string) error {
	// Create jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			semconv.ServiceNamespaceKey.String(namespace),
		)),
		trace.WithSampler(&CustomSampler{}),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(b3.New())
	return nil
}
