package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
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
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(input, fields...), pcs[0])
	_ = c.log.Handler().Handle(context.Background(), r)
}

func (c *logger) Error(input string, fields ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(input, fields...), pcs[0])
	_ = c.log.Handler().Handle(context.Background(), r)
}

func (c *logger) Debug(input string, fields ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(input, fields...), pcs[0])
	_ = c.log.Handler().Handle(context.Background(), r)
}
