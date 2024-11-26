package app

import (
	"github.com/satont/twir/apps/timers/internal/activity"
	bus_listener "github.com/satont/twir/apps/timers/internal/bus-listener"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	"github.com/satont/twir/apps/timers/internal/worker"
	"github.com/satont/twir/apps/timers/internal/workflow"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"timers",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "timers"}),
	fx.Provide(
		timers.NewGorm,
		activity.New,
		workflow.New,
		channels.NewGorm,
		streams.NewGorm,
		func(config cfg.Config) parser.ParserClient {
			return clients.NewParser(config.AppEnv)
		},
	),
	fx.Invoke(
		uptrace.NewFx("timers"),
		worker.New,
		bus_listener.New,
		func(l logger.Logger) {
			l.Info("Timers service started")
		},
	),
)
