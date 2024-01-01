package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	cfg "github.com/satont/twir/libs/config"
)

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
	Debug(input string, fields ...any)
	Warn(input string, fields ...any)
	WithComponent(name string) Logger
	GetSlog() *slog.Logger
}

type Log struct {
	log *slog.Logger

	service string
	sentry  *sentry.Client
}

type Opts struct {
	Env     string
	Service string

	Sentry *sentry.Client
	Level  slog.Level
}

func NewFx(opts Opts) func(config cfg.Config, sentry *sentry.Client) Logger {
	return func(config cfg.Config, sentry *sentry.Client) Logger {
		return New(
			Opts{
				Env:     config.AppEnv,
				Service: opts.Service,
				Sentry:  sentry,
			},
		)
	}
}

func New(opts Opts) Logger {
	level := opts.Level

	var zeroLogWriter io.Writer
	if opts.Env == "production" {
		zeroLogWriter = os.Stderr
	} else {
		zeroLogWriter = zerolog.ConsoleWriter{Out: os.Stderr}
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	slogzerolog.SourceKey = "source"
	slogzerolog.ErrorKeys = []string{"error", "err"}
	zerolog.ErrorStackFieldName = "stack"

	zeroLogLogger := zerolog.New(zeroLogWriter)

	log := slog.New(
		slogmulti.Fanout(
			slogzerolog.Option{
				Level:     level,
				Logger:    &zeroLogLogger,
				AddSource: true,
			}.NewZerologHandler(),
			slogsentry.Option{Level: slog.LevelError, AddSource: true}.NewSentryHandler(),
		),
	)

	if opts.Service != "" {
		log = log.With(slog.String("service", opts.Service))
	}

	return &Log{
		log:     log,
		sentry:  opts.Sentry,
		service: opts.Service,
	}
}

func (c *Log) handle(level slog.Level, input string, fields ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, input, pcs[0])
	for _, f := range fields {
		r.Add(f)
	}
	_ = c.log.Handler().Handle(context.Background(), r)
}

func (c *Log) Info(input string, fields ...any) {
	c.handle(slog.LevelInfo, input, fields...)
}

func (c *Log) Warn(input string, fields ...any) {
	c.handle(slog.LevelWarn, input, fields...)
}

func (c *Log) Error(input string, fields ...any) {
	c.handle(slog.LevelError, input, fields...)

}

func (c *Log) Debug(input string, fields ...any) {
	c.handle(slog.LevelDebug, input, fields...)
}

func (c *Log) WithComponent(name string) Logger {
	return &Log{
		log:     c.log.With(slog.String("component", name)),
		sentry:  c.sentry,
		service: c.service,
	}
}

func (c *Log) GetSlog() *slog.Logger {
	return c.log
}
