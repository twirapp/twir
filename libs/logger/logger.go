package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type Options struct {
	AppName     string
	Environment string
	Level       slog.Level
}

// New returns [*slog.Logger] with default handlers and given options.
//
// Default and additional handlers are assembled together and fan-outed through single handler.
func New(options Options, additionalHandlers ...slog.Handler) *slog.Logger {
	handlers := []slog.Handler{
		newZeroLogHandler(options.Level, options.Environment),
		// Sentry should handle only error log messages.
		newSentryHandler(slog.LevelError),
	}

	handlers = append(handlers, additionalHandlers...)

	logger := slog.New(
		slogmulti.Fanout(handlers...),
	)

	if options.AppName != "" {
		logger = logger.With(
			slog.String("app", options.AppName),
		)
	}

	return logger
}

// SetDefault creates [*slog.Logger] with New constructor and sets it as default (global) logger.
//
// See [slog.SetDefault] for more detailed description.
func SetDefault(options Options, additionalHandlers ...slog.Handler) {
	logger := New(options, additionalHandlers...)
	slog.SetDefault(logger)
}

// WithComponent returns a [*slog.Logger] that includes component attribute with the given name in each
// subsequent output operation.
func WithComponent(logger *slog.Logger, name string) *slog.Logger {
	return logger.With(
		slog.String("component", name),
	)
}

// Error returns a [slog.Attr] that represents error.
//
// In perfect world every log message that contains error output should be built with this helper.
func Error(err error) slog.Attr {
	if err == nil {
		return slog.Any("error", err)
	}

	return slog.String("error", err.Error())
}

func newSentryHandler(level slog.Level) slog.Handler {
	option := slogsentry.Option{
		Level:     level,
		AddSource: true,
	}

	return option.NewSentryHandler()
}

func newZeroLogHandler(level slog.Level, environment string) slog.Handler {
	var writer io.Writer

	switch environment {
	case "production":
		writer = os.Stderr
	default:
		writer = zerolog.ConsoleWriter{
			Out: os.Stderr,
		}
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	logger := zerolog.New(writer)
	option := slogzerolog.Option{
		Level:     level,
		Logger:    &logger,
		AddSource: true,
	}

	return option.NewZerologHandler()
}
