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

func NewFx(config cfg.Config) (*sentry.Client, error) {
	if config.SentryDsn == "" {
		return nil, nil
	}

	s, err := sentry.NewClient(
		sentry.ClientOptions{
			Dsn:              config.SentryDsn,
			AttachStacktrace: true,
		},
	)

	return s, err
}
