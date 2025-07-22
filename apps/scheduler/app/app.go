package app

import (
	bus_listener "github.com/satont/twir/apps/scheduler/internal/bus-listener"
	"github.com/satont/twir/apps/scheduler/internal/services"
	"github.com/satont/twir/apps/scheduler/internal/timers"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"

	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
)

const service = "scheduler"

var App = fx.Module(
	service,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: service}),
	fx.Provide(
		services.NewRoles,
		services.NewCommands,
		fx.Annotate(
			scheduledvipsrepositorypgx.NewFx,
			fx.As(new(scheduledvipsrepository.Repository)),
		),
	),
	fx.Invoke(
		uptrace.NewFx(service),
		bus_listener.New,
		timers.NewOnlineUsers,
		timers.NewStreams,
		timers.NewCommandsAndRoles,
		timers.NewWatched,
		timers.NewExpiredCommands,
		timers.NewScheduledVips,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)
