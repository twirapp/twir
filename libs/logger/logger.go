package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
}

type logger struct {
	log *slog.Logger
}

func New(env string) Logger {
	var log *slog.Logger

	switch env {
	case "development":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "production":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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
