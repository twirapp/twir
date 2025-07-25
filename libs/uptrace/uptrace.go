package uptrace

import (
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// nolint:ireturn
func New(config cfg.Config, service string) trace.Tracer {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(config.UptraceDsn),
		uptrace.WithServiceName(service),
		uptrace.WithDeploymentEnvironment(config.AppEnv),
		uptrace.WithTracingEnabled(true),
		uptrace.WithLoggingEnabled(true),
		uptrace.WithMetricsEnabled(true),
	)

	return otel.Tracer(service)
}

func NewFx(service string) func(config cfg.Config, lc fx.Lifecycle) trace.Tracer {
	return func(config cfg.Config, lc fx.Lifecycle) trace.Tracer {
		lc.Append(
			fx.Hook{
				OnStop: uptrace.Shutdown,
			},
		)

		return New(config, service)
	}
}
