package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
	Debug(input string, fields ...any)
}

type logger struct {
	log *slog.Logger
}

type Opts struct {
	Env     string
	Service string
}

func New(opts Opts) Logger {
	var log *slog.Logger

	switch opts.Env {
	case "development":
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug, AddSource: true,
				},
			),
		)
	case "production":
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true},
			),
		)
	default:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true},
			),
		)
	}

	if opts.Service != "" {
		log = log.With(slog.String("service", opts.Service))
	}

	return &logger{
		log: log,
	}
}

func (c *logger) Info(input string, fields ...any) {
	c.log.Info(input, fields...)
}

func (c *logger) Error(input string, fields ...any) {
	c.log.Error(input, fields...)
}

func (c *logger) Debug(input string, fields ...any) {
	c.log.Debug(input, fields...)
}
