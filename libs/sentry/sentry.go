package twirsentry

import (
	"log/slog"

	"github.com/getsentry/sentry-go"
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
		Debug:            false,
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
