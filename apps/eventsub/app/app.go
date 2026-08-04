package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/goforj/wire"
	"github.com/redis/go-redis/v9"
	buslistener "github.com/twirapp/twir/apps/eventsub/internal/bus-listener"
	"github.com/twirapp/twir/apps/eventsub/internal/handler"
	httpserver "github.com/twirapp/twir/apps/eventsub/internal/http"
	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	"github.com/twirapp/twir/apps/eventsub/internal/manager"
	eventplatforms "github.com/twirapp/twir/apps/eventsub/internal/platforms"
	usercreator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	"github.com/twirapp/twir/apps/eventsub/internal/vkvideo"
	"github.com/twirapp/twir/apps/eventsub/internal/webhook"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	channelsongrequestssettingscache "github.com/twirapp/twir/libs/cache/channel_song_requests_settings"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	channelsintegrationsseventvcache "github.com/twirapp/twir/libs/cache/channels_integrations_settings_seventv"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	genericcacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	grpcwebsockets "github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/integrations/vk"
	platformsregistry "github.com/twirapp/twir/libs/platforms"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	channelplatformsrepository "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsrepositorypgx "github.com/twirapp/twir/libs/repositories/channel_platforms/pgx"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelscommandsprefixrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelseventsrepository "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	channelsinfohistoryrepository "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsinfohistoryrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_info_history/datasource/postgres"
	channelsredemptionsrepository "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	channelsredemptionsrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_redemptions_history/datasources/clickhouse"
	commandsrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandsrepositorypgx "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	kickbotsrepositorypgx "github.com/twirapp/twir/libs/repositories/kick_bots/pgx"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypgx "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	twitchconduitsrepository "github.com/twirapp/twir/libs/repositories/twitch_conduits"
	twitchconduitsrepositorypgx "github.com/twirapp/twir/libs/repositories/twitch_conduits/datasource/postgres"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	usersstatsrepository "github.com/twirapp/twir/libs/repositories/users_stats"
	usersstatsrepositorypgx "github.com/twirapp/twir/libs/repositories/users_stats/datasources/postgres"
	userswithstatsrepository "github.com/twirapp/twir/libs/repositories/userswithstats"
	userswithstatsrepositorypgx "github.com/twirapp/twir/libs/repositories/userswithstats/datasource/postgres"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

const Service = "eventsub"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	channelplatformsrepositorypgx.NewFx,
	wire.Bind(new(channelplatformsrepository.Repository), new(*channelplatformsrepositorypgx.Pgx)),
	channelscommandsprefixrepositorypgx.NewFx,
	wire.Bind(
		new(channelscommandsprefixrepository.Repository),
		new(*channelscommandsprefixrepositorypgx.Pgx),
	),
	scheduledvipsrepositorypgx.NewFx,
	wire.Bind(new(scheduledvipsrepository.Repository), new(*scheduledvipsrepositorypgx.Pgx)),
	channelsinfohistoryrepositorypgx.NewFx,
	wire.Bind(
		new(channelsinfohistoryrepository.Repository),
		new(*channelsinfohistoryrepositorypgx.Pgx),
	),
	streamsrepositorypgx.NewFx,
	wire.Bind(new(streamsrepository.Repository), new(*streamsrepositorypgx.Pgx)),
	channelseventsrepositorypgx.NewFx,
	wire.Bind(new(channelseventsrepository.Repository), new(*channelseventsrepositorypgx.Pgx)),
	alertsrepositorypgx.NewFx,
	wire.Bind(new(alertsrepository.Repository), new(*alertsrepositorypgx.Pgx)),
	channelsredemptionsrepositoryclickhouse.NewFx,
	wire.Bind(
		new(channelsredemptionsrepository.Repository),
		new(*channelsredemptionsrepositoryclickhouse.Clickhouse),
	),
	twitchconduitsrepositorypgx.NewFx,
	wire.Bind(new(twitchconduitsrepository.Repository), new(*twitchconduitsrepositorypgx.Pgx)),
	commandsrepositorypgx.NewFx,
	wire.Bind(new(commandsrepository.Repository), new(*commandsrepositorypgx.Pgx)),
	userswithstatsrepositorypgx.NewFx,
	wire.Bind(new(userswithstatsrepository.Repository), new(*userswithstatsrepositorypgx.Pgx)),
	usersstatsrepositorypgx.NewFx,
	wire.Bind(new(usersstatsrepository.Repository), new(*usersstatsrepositorypgx.Pgx)),
	usersrepositorypgx.NewFx,
	wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
	kickbotsrepositorypgx.NewFx,
	wire.Bind(new(kickbotsrepository.Repository), new(*kickbotsrepositorypgx.Pgx)),
	channelcache.New,
	channelservice.NewChannelService,
	NewWebsocketClient,
	NewCommandsPrefixCache,
	wire.Struct(new(usercreator.Opts), "*"),
	usercreator.New,
	channelalertscache.New,
	commandscache.New,
	channelsongrequestssettingscache.New,
	channelsintegrationsseventvcache.New,
	wire.Struct(new(handler.Opts), "*"),
	handler.New,
	wire.Struct(new(manager.Opts), "*"),
	manager.NewManager,
	wire.Struct(new(kick.Opts), "*"),
	kick.New,
	NewTransportRegistry,
	wire.Struct(new(kick.HandlersOpts), "*"),
	kick.NewHandlers,
	httpserver.New,
	wire.Struct(new(kick.ResubscribeJobOpts), "*"),
	kick.NewResubscribeJob,
	wire.Struct(new(webhook.Opts), "*"),
	webhook.NewManager,
	wire.Struct(new(buslistener.Opts), "*"),
	buslistener.New,
	NewServerRunner,
	NewApplication,
)

type ServerRunner struct{}

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewWebsocketClient(config config.Config) grpcwebsockets.WebsocketClient {
	return clients.NewWebsocket(config.AppEnv)
}

func NewCommandsPrefixCache(
	repository channelscommandsprefixrepository.Repository,
	bus *buscore.Bus,
) *genericcacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix] {
	return channelscommandsprefixcache.New(repository, bus)
}

func NewTransportRegistry(
	config config.Config,
	kickTransport *kick.SubscriptionManager,
	logger *slog.Logger,
	redisClient *redis.Client,
	bus *buscore.Bus,
	userCreator *usercreator.UserCreatorService,
	channelsRepository channelsrepository.Repository,
	lc *lifecycle.Lifecycle,
) (*platformsregistry.Registry[eventplatforms.EventTransport], error) {
	return eventplatforms.NewVKVideoRegistry(
		config,
		kickTransport,
		func() (eventplatforms.EventTransport, error) {
			webSocketTokenClient, err := vk.NewWebSocketTokenClient(vk.VideoChatClientOpts{
				APIBaseURL: config.VKVideoDevAPIBaseURL,
			})
			if err != nil {
				return nil, fmt.Errorf("create VK Video WebSocket token client: %w", err)
			}

			return vkvideo.New(vkvideo.Opts{
				Logger:               logger,
				Redis:                redisClient,
				Bus:                  bus,
				UserCreator:          userCreator,
				WebSocketTokenClient: webSocketTokenClient,
				ChannelsRepo:         channelsRepository,
				Lc:                   lc,
				ProxyUrl:             config.VkProxyUrl,
			})
		},
	)
}

func NewServerRunner(server *httpserver.Server, lc *lifecycle.Lifecycle) *ServerRunner {
	lc.Append(lifecycle.Hook{
		OnStart: func(context.Context) error {
			return server.Start()
		},
		OnStop: server.Stop,
	})
	return &ServerRunner{}
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	_ *buslistener.BusListener,
	_ *ServerRunner,
	_ *webhook.Manager,
	_ *kick.ResubscribeJob,
) *Application {
	logger.Info("🚀 EventSub App started")
	return &Application{lifecycle: lifecycle}
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
