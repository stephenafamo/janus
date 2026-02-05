package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/stephenafamo/janus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type MetricsConfig struct {
	EnableOtelMetrics  bool          `env:"ENABLE_OTEL_METRICS"`
	MetricReadInterval time.Duration `env:"OTEL_METRICS_READ_INTERVAL,default=60s"`
}

func NewMeterProvider(ctx context.Context, config MetricsConfig, l Logger) (metric.MeterProvider, janus.StopFunc, error) {
	// Create resource.
	res, err := getOtelResource(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("creating resource: %w", err)
	}

	opts := []sdkmetric.Option{
		sdkmetric.WithResource(res),
	}

	if config.EnableOtelMetrics {
		l.Info(ctx, "Enabling OpenTelemetry Metrics Reader")

		metricExporter, err := otlpmetricgrpc.New(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("creating OTLP metric exporter: %w", err)
		}

		opts = append(opts,
			sdkmetric.WithReader(
				sdkmetric.NewPeriodicReader(
					metricExporter, sdkmetric.WithInterval(config.MetricReadInterval),
				),
			),
		)
	}

	meterProvider := sdkmetric.NewMeterProvider(opts...)

	// Register as global meter provider so that it can be used via otel.Meter
	// and accessed using otel.GetMeterProvider.
	// Most instrumentation libraries use the global meter provider as default.
	// If the global meter provider is not set then a no-op implementation
	// is used, which fails to generate data.
	otel.SetMeterProvider(meterProvider)

	return meterProvider, meterProvider.Shutdown, nil
}
