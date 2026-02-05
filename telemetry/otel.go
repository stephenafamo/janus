package telemetry

import (
	"context"
	"log/slog"
	"runtime/debug"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"go.opentelemetry.io/otel/trace"
)

func getOtelResource(ctx context.Context) (*resource.Resource, error) {
	options := []resource.Option{
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithProcessRuntimeDescription(),
		resource.WithTelemetrySDK(),
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		options = append(options, resource.WithAttributes(
			semconv.ServiceVersionKey.String(info.Main.Version),
		))
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // Propagates TraceParent
		propagation.Baggage{},      // Propagates Baggage
	))

	return resource.New(ctx, options...)
}

// ----------------------------------------------------------------------------
// SetAttrs does 3 things:
//
// 1. Sets slog attributes into the current span as span attributes.
// 2. Sets slog attributes into the current context as baggage members.
// 3. Stores slog attributes into the context for later retrieval by slog handlers.
//
// It returns the updated context.
// ----------------------------------------------------------------------------
func SetAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	spanAttrs := slogAttrToSpanAttr(attrs...)

	members := make([]baggage.Member, len(spanAttrs))
	for i, kv := range spanAttrs {
		memb, err := baggage.NewMember(string(kv.Key), kv.Value.Emit())
		if err != nil {
			slog.ErrorContext(
				ctx, "creating baggage member",
				slog.String("key", string(kv.Key)),
				slog.String("value", kv.Value.Emit()),
				slog.String("err", err.Error()),
			)
			continue
		}
		members[i] = memb
	}

	bag, err := baggage.New(append(baggage.FromContext(ctx).Members(), members...)...)
	if err != nil {
		slog.ErrorContext(ctx, "creating baggage", slog.String("err", err.Error()))
	} else {
		ctx = baggage.ContextWithBaggage(ctx, bag)
	}

	span := trace.SpanFromContext(ctx)
	span.SetAttributes(spanAttrs...)

	if existing, ok := ctx.Value(logAttrsKey).([]slog.Attr); ok {
		attrs = append(existing, attrs...)
	}

	return context.WithValue(ctx, logAttrsKey, attrs)
}
