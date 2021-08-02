package warehouse

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// InitializeJaeger creates a new trace provider instance and registers it as global trace provider.
func InitializeJaeger() (func(), bool) {
	jaegerUrl := os.Getenv("JAEGER_URL")
	if jaegerUrl == "" {
		log.Print("No Jaeger Tracing Url, will not report traces")
		return nil, false
	}

	otel.SetTextMapPropagator(b3.B3{})

	// Create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(jaegerUrl),
		jaeger.WithProcessFromEnv(),
	)
	if err != nil {
		log.Fatal(err)
	}
	return flush, true
}

// CreateSpan Create Tracing Span from context
func CreateSpan(ctx context.Context, tracerName, operationName string) (context.Context, oteltrace.Span) {
	tr := otel.Tracer(tracerName)
	return tr.Start(ctx, fmt.Sprintf("%s-%s", tracerName, operationName))
}

// GetSpan Gets the current Tracing Span from context
func GetSpan(ctx context.Context) oteltrace.Span {
	return oteltrace.SpanFromContext(ctx)
}
