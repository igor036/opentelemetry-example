package tracer

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	_zipikingURL = "http://localhost:9411/api/v2/spans"
	_appName     = "go-worker"
)

func InitTrace() (func(context.Context) error, error) {

	exporter, err := zipkin.New(_zipikingURL)

	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(_appName),
		)),
	)

	otel.SetTracerProvider(traceProvider)
	return traceProvider.Shutdown, nil
}

func Span(ctx context.Context, libary, name string) trace.Span {
	tr := otel.GetTracerProvider().Tracer("go-worder/client")
	_, span := tr.Start(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
	return span

}

func RunShutdown(shutdown func(context.Context) error, context context.Context) {
	if err := shutdown(context); err != nil {
		log.Fatal("failed to shutdown TracerProvider: %w", err)
	}
}

func HttpRequestSpanAttributes(span trace.Span, url, method, requestBody string) {
	span.SetAttributes(
		attribute.String("url", url),
		attribute.String("method", method),
		attribute.String("Request body", requestBody),
	)
}

func HttResponseSpanAttributes(span trace.Span, responseBody string, statusCode int) {
	span.SetAttributes(
		attribute.String("Response body", responseBody),
		attribute.Int("Response Status Code", statusCode),
	)
}

func HttErrorSpanAttributes(span trace.Span, err error) {
	span.SetAttributes(
		attribute.String("Error", err.Error()),
	)
}
