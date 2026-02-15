package twirsentry

import (
	"context"
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	cfg "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

func New(dsn, service string) (*sentry.Client, error) {
	if dsn == "" {
		slog.Warn("Sentry DSN is not set, Sentry will be disabled")
		return nil, nil
	}

	tags := map[string]string{}

	if service != "" {
		tags["service"] = service
	}

	o := sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Tags:             tags,
		Debug:            true,
		SendDefaultPII:   true,
		EnableLogs:       true,
		EnableTracing:    true,
	}

	s, err := sentry.NewClient(o)
	if err != nil {
		return nil, err
	}
	if err := sentry.Init(o); err != nil {
		return nil, err
	}

	return s, nil
}

type NewFxOpts struct {
	Service string
}

func NewFx(opts NewFxOpts) func(config cfg.Config, lc fx.Lifecycle) error {
	return func(config cfg.Config, lc fx.Lifecycle) error {
		_, err := New(config.SentryDsn, opts.Service)
		if err != nil {
			return err
		}

		lc.Append(
			fx.Hook{OnStop: func(ctx context.Context) error {
				sentry.Flush(2 * time.Second)

				return nil
			}},
		)

		return err
	}
}
