package logger

import (
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
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

type fxLogger struct {
	logger *slog.Logger
}

// NewFxLogger creates a new fx logger that wraps slog.Logger.
// It only logs error events from fx, all at Error level.
func NewFxLogger(logger *slog.Logger) fxevent.Logger {
	return &fxLogger{
		logger: logger.With(slog.String("component", "fx")),
	}
}

func (l *fxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.Error(
				"OnStart hook failed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.Any("error", e.Err),
			)
		}

	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.Error(
				"OnStop hook failed",
				slog.String("callee", e.FunctionName),
				slog.String("caller", e.CallerName),
				slog.Any("error", e.Err),
			)
		}

	case *fxevent.Supplied:
		if e.Err != nil {
			l.logger.Error(
				"Supply failed",
				slog.String("type", e.TypeName),
				slog.String("module", e.ModuleName),
				slog.Any("error", e.Err),
			)
		}

	case *fxevent.Provided:
		if e.Err != nil {
			l.logger.Error(
				"Provide failed",
				slog.String("constructor", e.ConstructorName),
				slog.String("module", e.ModuleName),
				slog.Any("error", e.Err),
			)
		}

	case *fxevent.Replaced:
		if e.Err != nil {
			l.logger.Error(
				"Replace failed",
				slog.String("module", e.ModuleName),
				slog.Any("error", e.Err),
			)
		}

	case *fxevent.Decorated:
		if e.Err != nil {
			l.logger.Error(
				"Decorate failed",
				slog.String("decorator", e.DecoratorName),
				slog.String("module", e.ModuleName),
				slog.Any("error", e.Err),
			)
		}

	// case *fxevent.:
	// 	if e.Err != nil {
	// 		l.logger.Error("Run failed",
	// 			slog.String("name", e.Name),
	// 			slog.String("kind", e.Kind),
	// 			slog.String("module", e.ModuleName),
	// 			slog.Any("error", e.Err),
	// 		)
	// 	}

	case *fxevent.Invoked:
		if e.Err != nil {
			l.logger.Error(
				"Invoke failed",
				slog.String("function", e.FunctionName),
				slog.String("module", e.ModuleName),
				slog.String("stack", e.Trace),
				slog.Any("error", e.Err),
			)
		}

	case *fxevent.Stopped:
		if e.Err != nil {
			l.logger.Error("Stop failed", slog.Any("error", e.Err))
		}

	case *fxevent.RollingBack:
		// Always log rollback - it's a critical failure
		l.logger.Error(
			"Start failed, rolling back",
			slog.Any("error", e.StartErr),
		)

	case *fxevent.RolledBack:
		if e.Err != nil {
			l.logger.Error("Rollback failed", slog.Any("error", e.Err))
		}

	case *fxevent.Started:
		if e.Err != nil {
			l.logger.Error("Start failed", slog.Any("error", e.Err))
		}

	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.logger.Error(
				"Logger initialization failed",
				slog.String("constructor", e.ConstructorName),
				slog.Any("error", e.Err),
			)
		}

		// All other events (success cases) are silently ignored
	}
}
