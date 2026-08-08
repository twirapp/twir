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
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
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
	"gorm.io/gorm"
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
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	commandsService *services.Commands,
	rolesService *services.Roles,
	bus *buscore.Bus,
	cfg config.Config,
	gorm *gorm.DB,
	streamsRepo streamsrepository.Repository,
	channelService *channelservice.ChannelService,
	commandsRepo commandsrepository.Repository,
	scheduledVipsRepo scheduledvipsrepository.Repository,
	usersRepo usersrepository.Repository,
) (*Application, error) {
	if err := buslistener.New(lifecycle, logger, commandsService, rolesService, bus); err != nil {
		return nil, fmt.Errorf("create scheduler listener: %w", err)
	}
	timers.NewOnlineUsers(lifecycle, logger, cfg, gorm, bus, streamsRepo)
	timers.NewStreams(lifecycle, cfg, logger, gorm, bus, streamsRepo, channelService)
	timers.NewCommandsAndRoles(lifecycle, logger, cfg, rolesService, commandsService, gorm)
	timers.NewWatched(lifecycle, logger, cfg, gorm, bus)
	timers.NewExpiredCommands(lifecycle, logger, cfg, gorm, bus, commandsRepo)
	timers.NewScheduledVips(lifecycle, cfg, logger, bus, scheduledVipsRepo, usersRepo, channelService)

	logger.Info("Started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
