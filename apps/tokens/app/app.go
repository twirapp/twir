package app

import (
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/go-redsync/redsync/v4"
	"github.com/goforj/wire"
	"github.com/twirapp/twir/apps/tokens/internal/bus_listener"
	redislock "github.com/twirapp/twir/apps/tokens/internal/redis"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
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
	"gorm.io/gorm"
)

const Service = "tokens"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
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
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	cfg config.Config,
	gorm *gorm.DB,
	redsync *redsync.Redsync,
	bus *buscore.Bus,
	kickBotsRepo kickbotsrepository.Repository,
	integrationsRepo integrationsrepository.Repository,
	channelIntegrationsRepo channelsintegrationsrepository.Repository,
	spotifyIntegrationsRepo channelsintegrationsspotifyrepository.Repository,
	tokensRepo tokensrepository.Repository,
	usersRepo usersrepository.Repository,
	vkVideoBotsRepo vkvideobotsrepository.Repository,
	youtubeBotsRepo youtubebotsrepository.Repository,
	trManager trm.Manager,
) (*Application, error) {
	if err := bus_listener.NewTokens(lifecycle, cfg, gorm, redsync, logger, bus, kickBotsRepo, integrationsRepo, channelIntegrationsRepo, spotifyIntegrationsRepo, tokensRepo, usersRepo, vkVideoBotsRepo, youtubeBotsRepo, trManager); err != nil {
		return nil, fmt.Errorf("create tokens listener: %w", err)
	}

	logger.Info("Started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
