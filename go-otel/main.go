package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0" // Use a stable semantic conventions version
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// initTracer initializes the OpenTelemetry TracerProvider.
func initTracer() *sdktrace.TracerProvider {
	// Create stdout exporter to be able to see the trace output.
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to create stdout exporter: %v", err)
	}

	// Create a resource describing your application.
	// This resource will be attached to all telemetry data.
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("go-otel-example-service"),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Create a new TracerProvider with the stdout exporter.
	// Use sdktrace.AlwaysSample to ensure all spans are recorded for demonstration.
	// In production, you might use a sampler like sdktrace.ParentBased or sdktrace.TraceIDRatioBased.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter), // Batcher for efficient export
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Register the TracerProvider globally. This allows easy access to the tracer.
	otel.SetTracerProvider(tp)
	// Set the context propagator to ensure trace context is propagated across network boundaries.
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}

// sayHello handles HTTP requests and creates spans.
func sayHello(w http.ResponseWriter, r *http.Request) {
	// Get the global tracer.
	tracer := otel.Tracer("my-app-tracer")

	// Start a custom span for this operation.
	// The parent context (r.Context()) will automatically link this span to the HTTP request span.
	ctx, span := tracer.Start(r.Context(), "sayHello-operation")
	defer span.End() // Ensure the span is ended when the function exits.

	// Add an attribute to the custom span.
	span.SetAttributes(attribute.String("request.path", r.URL.Path))

	log.Printf("Handling request for: %s", r.URL.Path)

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	// Start a child span for a sub-operation.
	_, subSpan := tracer.Start(ctx, "simulate-db-call")
	defer subSpan.End()
	subSpan.SetAttributes(attribute.Int("db.rows_affected", 10))
	time.Sleep(20 * time.Millisecond)

	fmt.Fprintf(w, "Hello, OpenTelemetry Go!")
}

func main() {
	// Initialize the TracerProvider. Defer its shutdown to flush any remaining spans.
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Use otelhttp.NewHandler to automatically instrument incoming HTTP requests.
	// This will create a span for each incoming request.
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(sayHello), "hello-endpoint"))

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}