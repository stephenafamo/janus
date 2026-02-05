package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/go-logr/logr"
	slogstrict "github.com/go-swiss/slog-strict"
	slogotel "github.com/remychantenay/slog-otel"
	slogmulti "github.com/samber/slog-multi"
	"github.com/stephenafamo/janus"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/processors/minsev"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	logsdk "go.opentelemetry.io/otel/sdk/log"
)

type ctxKey string

const logAttrsKey ctxKey = "log_attrs"

type Logger slogstrict.Logger

type LogConfig struct {
	Level             string `env:"LOG_LEVEL,default=info"`
	EnableOtelLogging bool   `env:"ENABLE_OTEL_LOGGING"`
}

type handlerSetupFunc = func(ctx context.Context, level slog.Level) (slog.Handler, janus.StopFunc, error)

// NewLogger adds openTelemetry logging support if enabled via environment variable. It returns a slogstrict.Logger and a cleanup function to flush logs on shutdown.
// Because this sets the default slog logger, it MUST NOT be called by passing the
// bare default slog logger as the handler, or it will deadlock.
// i.e. DO NOT do this:
//
//	NewLogger(ctx, cfg, func(..){return slog.Default().Handler(), nil, nil})
//
// Do this instead:
//
//	NewLogger(ctx, cfg, func(..){return slog.NewTextHandler(os.Stdout, nil), nil, nil})
func NewLogger(ctx context.Context, cfg LogConfig, base handlerSetupFunc) (slogstrict.Logger, janus.StopFunc, error) {
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, nil, fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	cleanup := janus.NoopStopFunc

	handler, baseCleanup, err := base(ctx, level)
	if err != nil {
		return nil, nil, fmt.Errorf("creating base log handler: %w", err)
	}

	if cfg.EnableOtelLogging {
		otelHandler, otelCleanup, err := getOtelHandler(ctx, level)
		if err != nil {
			return nil, baseCleanup, fmt.Errorf("creating log handler: %w", err)
		}

		// Create a handler that sends log entries to both the terminal and sentry
		// Replace with slog.MultiHandler in Go 1.26
		handler = slogmulti.Fanout(handler, otelHandler)
		cleanup = janus.CombineStopFuncs(baseCleanup, otelCleanup)
	}

	// Add OpenTelemetry handler wrapper
	// This adds trace and span IDs to log records automatically
	handler = slogotel.OtelHandler{Next: handler, NoBaggage: true}

	// Add Context handler wrapper
	// This adds attributes from context to log records automatically
	handler = newContextHandler(handler)

	sl := slog.New(handler)
	slog.SetDefault(sl)

	return slogstrict.FromSlog(sl), cleanup, nil
}

// ----------------------------------------------------------------------------
// Context Handler Middleware
// ----------------------------------------------------------------------------
type contextHandler struct {
	handler slog.Handler
}

// newContextHandler returns a new ContextHandler wrapping the given handler.
func newContextHandler(h slog.Handler) *contextHandler {
	// Optimization: If the handler is already a ContextHandler, just return it.
	if lh, ok := h.(*contextHandler); ok {
		return lh
	}
	return &contextHandler{handler: h}
}

// Enabled delegates to the underlying handler.
func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements the middleware logic.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(logAttrsKey).([]slog.Attr); ok {
		r.AddAttrs(attrs...)
	}

	return h.handler.Handle(ctx, r)
}

// WithAttrs ensures our ContextHandler wrapper persists when attributes are added.
func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newContextHandler(h.handler.WithAttrs(attrs))
}

// WithGroup ensures our ContextHandler wrapper persists when a group is added.
func (h *contextHandler) WithGroup(name string) slog.Handler {
	return newContextHandler(h.handler.WithGroup(name))
}

// ----------------------------------------------------------------------------
// Otel Logging
// ----------------------------------------------------------------------------
var _ minsev.Severitier = otelSeverity{}

type otelSeverity struct{ severity log.Severity }

// Severity implements minsev.Severitier.
func (o otelSeverity) Severity() log.Severity {
	return o.severity
}

func getOtelHandler(ctx context.Context, level slog.Level) (slog.Handler, janus.StopFunc, error) {
	resource, err := getOtelResource(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("creating resource: %w", err)
	}

	severity := log.SeverityDebug + log.Severity(level-slog.LevelDebug)

	// Create an exporter that will emit log records.
	// E.g. use go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp
	// to send logs using OTLP over HTTP:
	exporter, err := otlploggrpc.New(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("creating OTLP log exporter: %w", err)
	}

	// Wrap the processor so that it filters by severity level defined
	// via environmental variable.
	processor := minsev.NewLogProcessor(
		logsdk.NewBatchProcessor(exporter),
		otelSeverity{severity: severity},
	)

	// Create a logger provider.
	// You can pass this instance directly when creating a log bridge.
	provider := logsdk.NewLoggerProvider(
		logsdk.WithResource(resource),
		logsdk.WithProcessor(processor),
	)

	handler := otelslog.NewHandler("root", otelslog.WithLoggerProvider(provider))

	// Register as global logger provider so that it can be used via global.Meter
	// and accessed using global.GetLoggerProvider.
	// Most log bridges use the global logger provider as default.
	// If the global logger provider is not set then a no-op implementation
	// is used, which fails to generate data.
	global.SetLoggerProvider(provider)

	// Set the slog handler as the global otel logger.
	otel.SetLogger(logr.FromSlogHandler(handler))

	// Set the global otel error handler to log errors using slog.
	errorlogger := slog.New(handler).With(slog.String("by", "otel error handler"))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		errorlogger.Error(err.Error())
	}))

	return handler, provider.Shutdown, nil
}

func slogAttrToSpanAttr(attr ...slog.Attr) []attribute.KeyValue {
	spanAttrs := make([]attribute.KeyValue, 0, len(attr))
	for _, attr := range attr {
		switch attr.Value.Kind() {
		case slog.KindGroup:
			// Skip groups for now.
			spanAttrs = append(spanAttrs, slogAttrToSpanAttr(attr.Value.Group()...)...)

		default:
			spanAttrs = append(spanAttrs, attribute.KeyValue{
				Key:   attribute.Key(attr.Key),
				Value: slogValueToSpanValue(attr.Value),
			})
		}
	}

	return spanAttrs
}

func slogValueToSpanValue(value slog.Value) attribute.Value {
	/*
		KindAny Kind = iota
			KindBool
			KindDuration
			KindFloat64
			KindInt64
			KindString
			KindTime
			KindUint64
			KindGroup
			KindLogValuer
	*/
	switch value.Kind() {
	case slog.KindBool:
		return attribute.BoolValue(value.Bool())

	case slog.KindDuration:
		return attribute.Int64Value(int64(value.Duration()))

	case slog.KindFloat64:
		return attribute.Float64Value(value.Float64())

	case slog.KindInt64:
		return attribute.Int64Value(value.Int64())

	case slog.KindString:
		return attribute.StringValue(value.String())

	case slog.KindTime:
		return attribute.StringValue(value.Time().Format(time.RFC3339Nano))

	case slog.KindUint64:
		return attribute.Int64Value(int64(value.Uint64()))

	case slog.KindLogValuer:
		return slogValueToSpanValue(value.LogValuer().LogValue())

	default:
		return attribute.StringValue(fmt.Sprintf("%v", value.Any()))
	}
}
