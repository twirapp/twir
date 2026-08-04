package app

import (
	"fmt"
	"log/slog"

	"github.com/goforj/wire"
	buslistener "github.com/twirapp/twir/apps/timers/internal/bus-listener"
	"github.com/twirapp/twir/apps/timers/internal/manager"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypgx "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

const Service = "timers"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	timersrepositorypgx.NewFx,
	wire.Bind(new(timersrepository.Repository), new(*timersrepositorypgx.Pgx)),
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	streamsrepositorypgx.NewFx,
	wire.Bind(new(streamsrepository.Repository), new(*streamsrepositorypgx.Pgx)),
	channelservice.NewChannelService,
	channelcache.New,
	manager.New,
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	bus *buscore.Bus,
	managerService *manager.Manager,
) (*Application, error) {
	if err := buslistener.New(lifecycle, logger, bus, managerService); err != nil {
		return nil, fmt.Errorf("create timers listener: %w", err)
	}

	logger.Info("🚀 Timers service started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
