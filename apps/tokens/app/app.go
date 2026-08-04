package app

import (
	"fmt"
	"log/slog"

	"github.com/goforj/wire"
	"github.com/twirapp/twir/apps/tokens/internal/bus_listener"
	redislock "github.com/twirapp/twir/apps/tokens/internal/redis"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	channelsintegrationsrepository "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_integrations/datasource/postgres"
	channelsintegrationsspotifyrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifyrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	integrationsrepository "github.com/twirapp/twir/libs/repositories/integrations"
	integrationsrepositorypgx "github.com/twirapp/twir/libs/repositories/integrations/datasource/postgres"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	kickbotsrepositorypgx "github.com/twirapp/twir/libs/repositories/kick_bots/pgx"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokensrepositorypgx "github.com/twirapp/twir/libs/repositories/tokens/datasources/postgres"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	vkvideobotsrepository "github.com/twirapp/twir/libs/repositories/vk_video_bots"
	vkvideobotsrepositorypgx "github.com/twirapp/twir/libs/repositories/vk_video_bots/datasource/postgres"
	youtubebotsrepository "github.com/twirapp/twir/libs/repositories/youtube_bots"
	youtubebotsrepositorypgx "github.com/twirapp/twir/libs/repositories/youtube_bots/datasource/postgres"
)

const Service = "tokens"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.NewBase,
	wire.FieldsOf(
		new(baseapp.Base),
		"Lifecycle",
		"Config",
		"Gorm",
		"Redis",
		"Logger",
		"Bus",
		"PgxPool",
		"TrManager",
	),
	tokensrepositorypgx.NewFx,
	wire.Bind(new(tokensrepository.Repository), new(*tokensrepositorypgx.Pgx)),
	usersrepositorypgx.NewFx,
	wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
	kickbotsrepositorypgx.NewFx,
	wire.Bind(new(kickbotsrepository.Repository), new(*kickbotsrepositorypgx.Pgx)),
	vkvideobotsrepositorypgx.NewFx,
	wire.Bind(new(vkvideobotsrepository.Repository), new(*vkvideobotsrepositorypgx.Pgx)),
	youtubebotsrepositorypgx.NewFx,
	wire.Bind(new(youtubebotsrepository.Repository), new(*youtubebotsrepositorypgx.Pgx)),
	channelsintegrationsrepositorypgx.NewFx,
	wire.Bind(
		new(channelsintegrationsrepository.Repository),
		new(*channelsintegrationsrepositorypgx.Pgx),
	),
	channelsintegrationsspotifyrepositorypgx.NewFx,
	wire.Bind(
		new(channelsintegrationsspotifyrepository.Repository),
		new(*channelsintegrationsspotifyrepositorypgx.Pgx),
	),
	integrationsrepositorypgx.NewFx,
	wire.Bind(new(integrationsrepository.Repository), new(*integrationsrepositorypgx.Pgx)),
	redislock.NewRedisLock,
	wire.Struct(new(bus_listener.Opts), "*"),
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	listenerOpts bus_listener.Opts,
) (*Application, error) {
	if err := bus_listener.NewTokens(listenerOpts); err != nil {
		return nil, fmt.Errorf("create tokens listener: %w", err)
	}

	logger.Info("Started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
