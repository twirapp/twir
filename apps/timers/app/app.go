package app

import (
	"github.com/twirapp/twir/apps/timers/internal/activity"
	bus_listener "github.com/twirapp/twir/apps/timers/internal/bus-listener"
	"github.com/twirapp/twir/apps/timers/internal/repositories/channels"
	"github.com/twirapp/twir/apps/timers/internal/repositories/streams"
	"github.com/twirapp/twir/apps/timers/internal/worker"
	"github.com/twirapp/twir/apps/timers/internal/workflow"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"timers",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "timers"}),
	fx.Provide(
		fx.Annotate(
			timersrepositorypgx.NewFx,
			fx.As(new(timersrepository.Repository)),
		),
		activity.New,
		workflow.New,
		channels.NewGorm,
		streams.NewGorm,
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
