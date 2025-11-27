package logger

import (
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type FxOptions struct {
	fx.In

	AdditionalHandlers []slog.Handler `group:"slog-handlers"`
}

// NewFx returns [*slog.Logger] via New with the given options and additional handlers from fx constructors.
func NewFx(options Options) func(fxOptions FxOptions) *slog.Logger {
	return func(fxOptions FxOptions) *slog.Logger {
		return New(
			Options{
				AppName:     options.AppName,
				Environment: options.Environment,
				Level:       options.Level,
			},
			fxOptions.AdditionalHandlers...,
		)
	}
}

// FxOnlyErrorsLoggerOption returns the [fx.WithLogger] fx option with [zap.Logger] with development mode
// and zap.ErrorLevel level set.
func FxOnlyErrorsLoggerOption() fx.Option {
	return fx.WithLogger(
		func() fxevent.Logger {
			config := zap.NewDevelopmentConfig()
			config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)

			return &fxevent.ZapLogger{
				Logger: zap.Must(config.Build()),
			}
		},
	)
}
