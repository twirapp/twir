package main

import (
	"log/slog"

	"github.com/satont/twir/apps/timers/internal/activity"
	"github.com/satont/twir/apps/timers/internal/gorm"
	"github.com/satont/twir/apps/timers/internal/grpc_server"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	"github.com/satont/twir/apps/timers/internal/worker"
	"github.com/satont/twir/apps/timers/internal/workflow"
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
			logger.NewFx(logger.Opts{Level: slog.LevelInfo, Service: "timers"}),
			gorm.New,
			timers.NewGorm,
			activity.New,
			workflow.New,
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
			worker.New,
			grpc_server.New,
			func(l logger.Logger) {
				l.Info("Timers service started")
			},
		),
	).Run()
}
