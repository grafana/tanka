package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("tanka")
}

func Start(ctx context.Context, name string, attributes ...attribute.KeyValue) (context.Context, trace.Span) {
	newCtx, span := tracer.Start(ctx, name)
	span.SetAttributes(attributes...)
	return newCtx, span
}

func InstallExportPipeline(ctx context.Context) (func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OTLP exporter: %w", err)
	}

	traceResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("tanka"),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(traceResource),
	)
	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}
