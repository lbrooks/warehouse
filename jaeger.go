package warehouse

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type jaegerConfig struct {
	url     string
	process jaeger.Process
}

func getJaegerConfig(service string) jaegerConfig {
	jaegerEndpoint := os.Getenv("JAEGER_URL")
	if jaegerEndpoint == "" {
		panic("JAEGER_URL is undefined")
	}

	host, err := os.Hostname()
	if err != nil {
		panic("Could not determine host")
	}

	return jaegerConfig{
		url: jaegerEndpoint,
		process: jaeger.Process{
			ServiceName: service,
			Tags: []label.KeyValue{
				label.String("host", host),
			},
		},
	}
}

// InitializeJaeger creates a new trace provider instance and registers it as global trace provider.
func InitializeJaeger(service string) func() {
	config := getJaegerConfig(service)

	otel.SetTextMapPropagator(b3.B3{})

	// Create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(config.url),
		jaeger.WithProcess(config.process),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return flush
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
