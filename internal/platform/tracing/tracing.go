package tracing

import (
	"context"
	"fmt"
	"github.com/alexperezortuno/go-url-shortner/internal/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StartSpan initializes a span for a given handler and updates the request context.
func StartSpan(ctx *gin.Context, serviceName string, spanName string) {
	tracer := otel.Tracer(spanName)
	spanCtx, span := tracer.Start(ctx.Request.Context(), serviceName)
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("http.method", ctx.Request.Method),
		attribute.String("http.route", ctx.FullPath()),
	)

	// Replace the context in the request with the span's context
	ctx.Request = ctx.Request.WithContext(spanCtx)
}

func InitTracer(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
	retryPolicy := `{
        "methodConfig": [{
            "name": [{"service": "opentelemetry.proto.collector.trace.v1.TraceService"}],
            "waitForReady": true,
            "retryPolicy": {
                "MaxAttempts": 5,
                "InitialBackoff": "1s",
                "MaxBackoff": "5s",
                "BackoffMultiplier": 2.0,
                "RetryableStatusCodes": ["UNAVAILABLE"]
            }
        }]
    }`

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(cfg.OtelExporterEndpoint),
		otlptracegrpc.WithDialOption(
			grpc.WithDefaultServiceConfig(retryPolicy),
		),
		otlptracegrpc.WithDialOption(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)

	if err != nil {
		return nil, err
	}

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global propagator and tracer provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Función de shutdown para limpieza
	shutdown := func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
		return nil
	}

	return shutdown, nil
}
