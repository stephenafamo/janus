package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/grafana/pyroscope-go"
	"github.com/stephenafamo/janus"
)

type ProfilingConfig struct {
	ServiceName   string `env:"OTEL_SERVICE_NAME"`
	ServerAddress string `env:"PYROSCOPE_SERVER_ADDRESS"`
}

func StartProfiling(ctx context.Context, config ProfilingConfig, l Logger) (janus.StopFunc, error) {
	if config.ServiceName == "" {
		l.Info(ctx, "Profiling is disabled because OTEL_SERVICE_NAME is not set")
		return janus.NoopStopFunc, nil
	}

	if config.ServerAddress == "" {
		l.Info(ctx, "Profiling is disabled because PYROSCOPE_SERVER_ADDRESS is not set")
		return janus.NoopStopFunc, nil
	}

	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: config.ServiceName,
		ServerAddress:   config.ServerAddress,

		// you can disable logging by setting this to nil
		Logger: pyroscopeLogger{l: l.ToSlog()},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start pyroscope profiler: %w", err)
	}

	return janus.NoContextStopFunc(profiler.Stop), nil
}

var _ pyroscope.Logger = &pyroscopeLogger{}

type pyroscopeLogger struct{ l *slog.Logger }

// Debugf implements pyroscope.Logger.
func (p pyroscopeLogger) Debugf(format string, args ...any) {
	p.l.Debug(fmt.Sprintf(format, args...))
}

// Errorf implements pyroscope.Logger.
func (p pyroscopeLogger) Errorf(format string, args ...any) {
	p.l.Error(fmt.Sprintf(format, args...))
}

// Infof implements pyroscope.Logger.
func (p pyroscopeLogger) Infof(format string, args ...any) {
	p.l.Info(fmt.Sprintf(format, args...))
}
