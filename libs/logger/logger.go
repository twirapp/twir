package logger

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	slogmulti "github.com/samber/slog-multi"
	slogsentry "github.com/samber/slog-sentry/v2"
	slogzerolog "github.com/samber/slog-zerolog/v2"
	"github.com/twirapp/twir/libs/logger/audit"
	"github.com/twirapp/twir/libs/logger/levels"
	"go.uber.org/fx"
)

type Logger interface {
	Info(input string, fields ...any)
	Error(input string, fields ...any)
	Debug(input string, fields ...any)
	Warn(input string, fields ...any)
	Audit(input string, fields audit.Fields)

	WithComponent(name string) Logger
	With(args ...any) Logger
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

	AdditionalSlogHandlers []slog.Handler
}

type FxOpts struct {
	fx.In

	Sentry                 *sentry.Client
	AdditionalSlogHandlers []slog.Handler `group:"slog-handlers"`
}

func NewFx(opts Opts) func(
	fxOpts FxOpts,
) Logger {
	return func(fxOpts FxOpts) Logger {
		return New(
			Opts{
				Env:                    opts.Env,
				Service:                opts.Service,
				Sentry:                 fxOpts.Sentry,
				AdditionalSlogHandlers: fxOpts.AdditionalSlogHandlers,
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

	handlers := []slog.Handler{
		slogzerolog.Option{
			Level:     level,
			Logger:    &zeroLogLogger,
			AddSource: true,
		}.NewZerologHandler(),
		slogsentry.Option{Level: slog.LevelError, AddSource: true}.NewSentryHandler(),
	}

	for _, h := range opts.AdditionalSlogHandlers {
		handlers = append(handlers, h)
	}

	log := slog.New(
		slogmulti.Fanout(handlers...),
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

func (c *Log) handle(ctx context.Context, level slog.Level, input string, fields ...any) {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, input, pcs[0])
	for _, f := range fields {
		r.Add(f)
	}

	_ = c.log.Handler().Handle(ctx, r)
}

func (c *Log) Info(input string, fields ...any) {
	c.handle(context.Background(), slog.LevelInfo, input, fields...)
}

func (c *Log) Warn(input string, fields ...any) {
	c.handle(context.Background(), slog.LevelWarn, input, fields...)
}

func (c *Log) Error(input string, fields ...any) {
	c.handle(context.Background(), slog.LevelError, input, fields...)
}

func (c *Log) Debug(input string, fields ...any) {
	c.handle(context.Background(), slog.LevelDebug, input, fields...)
}

func (c *Log) Audit(input string, fields audit.Fields) {
	ctx := context.WithValue(context.Background(), audit.AuditFieldsContextKey{}, fields)

	fieldsBytes, _ := json.Marshal(fields)
	fieldsMap := make(map[string]any)
	_ = json.Unmarshal(fieldsBytes, &fieldsMap)

	var castedFiles []any
	for k, v := range fieldsMap {
		castedFiles = append(castedFiles, slog.Any(strings.ToLower(k), v))
	}

	c.handle(ctx, levels.LevelAudit, input, castedFiles...)
}

func (c *Log) WithComponent(name string) Logger {
	return &Log{
		log:     c.log.With(slog.String("component", name)),
		sentry:  c.sentry,
		service: c.service,
	}
}

func (c *Log) With(args ...any) Logger {
	return &Log{
		log:     c.log.With(args...),
		sentry:  c.sentry,
		service: c.service,
	}
}

func (c *Log) GetSlog() *slog.Logger {
	return c.log
}
