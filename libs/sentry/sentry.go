package twirsentry

import (
	"github.com/getsentry/sentry-go"
	cfg "github.com/twirapp/twir/libs/config"
)

func New(dsn string) (*sentry.Client, error) {
	if dsn == "" {
		return nil, nil
	}

	opts := sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
	}

	s, err := sentry.NewClient(opts)

	sentry.Init(opts)

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

		o := sentry.ClientOptions{
			Dsn:              config.SentryDsn,
			AttachStacktrace: true,
			Tags:             tags,
			Debug:            false,
		}

		s, err := sentry.NewClient(o)
		sentry.Init(o)

		return s, err
	}
}
