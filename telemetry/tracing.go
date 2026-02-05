package telemetry

import (
	"context"
	"fmt"

	otelpyroscope "github.com/grafana/otel-profiling-go"
	"github.com/stephenafamo/janus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type TracingConfig struct {
	EnableOtelTracing bool    `env:"ENABLE_OTEL_TRACING"`
	TraceSampleRate   float64 `env:"OTEL_TRACES_SAMPLE_RATE,default=0.1"`
}

func NewTracerProvider(ctx context.Context, config TracingConfig, l Logger) (trace.TracerProvider, janus.StopFunc, error) {
	resource, err := getOtelResource(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("creating resource: %w", err)
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(config.TraceSampleRate))),
	}

	if config.EnableOtelTracing {
		l.Info(ctx, "Enabling OpenTelemetry Trace Exporter")

		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("creating OTLP trace exporter: %w", err)
		}
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	sdkProvider := sdktrace.NewTracerProvider(opts...)
	provider := otelpyroscope.NewTracerProvider(sdkProvider)

	otel.SetTracerProvider(provider)

	return provider, sdkProvider.Shutdown, nil
}
