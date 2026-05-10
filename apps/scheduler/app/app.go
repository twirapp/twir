package app

import (
	"log/slog"

	bus_listener "github.com/twirapp/twir/apps/scheduler/internal/bus-listener"
	"github.com/twirapp/twir/apps/scheduler/internal/services"
	"github.com/twirapp/twir/apps/scheduler/internal/timers"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/otel"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandswithgroupsandresponsespostgres "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	"go.uber.org/fx"

	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
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
		fx.Annotate(
			commandswithgroupsandresponsespostgres.NewFx,
			fx.As(new(commandswithgroupsandresponsesrepository.Repository)),
		),
		fx.Annotate(
			usersrepositorypgx.NewFx,
			fx.As(new(usersrepository.Repository)),
		),
		fx.Annotate(
			channelsrepositorypgx.NewFx,
			fx.As(new(channelsrepository.Repository)),
		),
	),
	fx.Invoke(
		otel.NewFx(service),
		bus_listener.New,
		timers.NewOnlineUsers,
		timers.NewStreams,
		timers.NewCommandsAndRoles,
		timers.NewWatched,
		timers.NewExpiredCommands,
		timers.NewScheduledVips,
		func(l *slog.Logger) {
			l.Info("Started")
		},
	),
)
