// Пример настройки OTLP экспортера в Go
package metrics

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	// "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitOTLP инициализирует OpenTelemetry с заданными параметрами
func InitOTLP(serviceName, otlpHost, otlpPort string) (*trace.TracerProvider, *metric.MeterProvider, error) {
	// Формируем URL для OTLP экспортеров
	otlpEndpoint := fmt.Sprintf("%s:%s", otlpHost, otlpPort)

	// Создаем OTLP экспортера для трассировки
	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(otlpEndpoint),
	)
	if err != nil {
		return nil, nil, err
	}

	// Создаем OTLP экспортера для метрик
	metricExporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint(otlpEndpoint),
	)
	if err != nil {
		return nil, nil, err
	}

	// Создаем ресурс с атрибутами
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	)

	// Создаем TracerProvider для трассировки
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)

	// Создаем MeterProvider для метрик
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(resource),
	)
	otel.SetMeterProvider(meterProvider)

	return tracerProvider, meterProvider, nil
}
