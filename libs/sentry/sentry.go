package twirsentry

import (
	"github.com/getsentry/sentry-go"
	cfg "github.com/satont/twir/libs/config"
)

func New(dsn string) (*sentry.Client, error) {
	if dsn == "" {
		return nil, nil
	}

	s, err := sentry.NewClient(
		sentry.ClientOptions{
			Dsn:              dsn,
			AttachStacktrace: true,
		},
	)

	return s, err
}

type NewFxOpts struct {
	Service string
}

func NewFx(opts NewFxOpts) func(config cfg.Config) (*sentry.Client, error) {
	return func(config cfg.Config) (*sentry.Client, error) {
		if config.SentryDsn == "" {
			return nil, nil
		}

		tags := map[string]string{}

		if opts.Service != "" {
			tags["service"] = opts.Service
		}

		s, err := sentry.NewClient(
			sentry.ClientOptions{
				Dsn:              config.SentryDsn,
				AttachStacktrace: true,
				Tags:             tags,
			},
		)

		return s, err
	}
}
