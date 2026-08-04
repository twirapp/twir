package app

import (
	"fmt"
	"log/slog"

	"github.com/goforj/wire"
	buslistener "github.com/twirapp/twir/apps/scheduler/internal/bus-listener"
	"github.com/twirapp/twir/apps/scheduler/internal/services"
	"github.com/twirapp/twir/apps/scheduler/internal/timers"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	commandsrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandsrepositorypgx "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypgx "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

const Service = "scheduler"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	scheduledvipsrepositorypgx.NewFx,
	wire.Bind(new(scheduledvipsrepository.Repository), new(*scheduledvipsrepositorypgx.Pgx)),
	commandsrepositorypgx.NewFx,
	wire.Bind(new(commandsrepository.Repository), new(*commandsrepositorypgx.Pgx)),
	usersrepositorypgx.NewFx,
	wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	streamsrepositorypgx.NewFx,
	wire.Bind(new(streamsrepository.Repository), new(*streamsrepositorypgx.Pgx)),
	channelservice.NewChannelService,
	services.NewRoles,
	services.NewCommands,
	wire.Struct(new(buslistener.Opts), "*"),
	wire.Struct(new(timers.OnlineUsersOpts), "*"),
	wire.Struct(new(timers.StreamOpts), "*"),
	wire.Struct(new(timers.CommandsAndRolesOpts), "*"),
	wire.Struct(new(timers.WatchedOpts), "*"),
	wire.Struct(new(timers.ExpiredCommandsOpts), "*"),
	wire.Struct(new(timers.ScheduledVipsOpts), "*"),
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	listenerOpts buslistener.Opts,
	onlineUsersOpts timers.OnlineUsersOpts,
	streamsOpts timers.StreamOpts,
	commandsAndRolesOpts timers.CommandsAndRolesOpts,
	watchedOpts timers.WatchedOpts,
	expiredCommandsOpts timers.ExpiredCommandsOpts,
	scheduledVipsOpts timers.ScheduledVipsOpts,
) (*Application, error) {
	if err := buslistener.New(listenerOpts); err != nil {
		return nil, fmt.Errorf("create scheduler listener: %w", err)
	}
	timers.NewOnlineUsers(onlineUsersOpts)
	timers.NewStreams(streamsOpts)
	timers.NewCommandsAndRoles(commandsAndRolesOpts)
	timers.NewWatched(watchedOpts)
	timers.NewExpiredCommands(expiredCommandsOpts)
	timers.NewScheduledVips(scheduledVipsOpts)

	logger.Info("Started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
