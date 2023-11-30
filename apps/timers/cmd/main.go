package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/satont/twir/apps/timers/internal/gorm"
	"github.com/satont/twir/apps/timers/internal/grpc_server"
	"github.com/satont/twir/apps/timers/internal/queue"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/logger"
	sentryInternal "github.com/satont/twir/libs/sentry"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			cfg.NewFx,
			sentryInternal.NewFx(sentryInternal.NewFxOpts{Service: "timers"}),
			func(config cfg.Config, s *sentry.Client) logger.Logger {
				return logger.New(
					logger.Opts{
						Env:     config.AppEnv,
						Service: "timers",
						Sentry:  s,
					},
				)
			},
			gorm.New,
			queue.New,
			timers.NewGorm,
			channels.NewGorm,
			streams.NewGorm,
			func(config cfg.Config) parser.ParserClient {
				return clients.NewParser(config.AppEnv)
			},
			func(config cfg.Config) bots.BotsClient {
				return clients.NewBots(config.AppEnv)
			},
		),
		fx.NopLogger,
		fx.Invoke(
			// queue.New,
			grpc_server.New,
		),
	).Run()
}
