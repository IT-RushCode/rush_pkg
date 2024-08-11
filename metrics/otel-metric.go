package metrics

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitTracer инициализирует трайсер.
func InitTracer(serviceName string) (*trace.TracerProvider, error) {
	tp := trace.NewTracerProvider(
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// ShutdownTracer завершает работу трайсера.
func ShutdownTracer(tp *trace.TracerProvider) {
	if err := tp.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown tracer provider: %v", err)
	}
}
