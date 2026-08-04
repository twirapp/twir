package app

import (
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"

	"github.com/goforj/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	bus_listener "github.com/twirapp/twir/apps/bots/internal/bus-listener"
	discordbushandler "github.com/twirapp/twir/apps/bots/internal/discord/bus_handler"
	"github.com/twirapp/twir/apps/bots/internal/discord/discord_go"
	"github.com/twirapp/twir/apps/bots/internal/discord/messages_updater"
	notificationssync "github.com/twirapp/twir/apps/bots/internal/discord/notifications_sync"
	"github.com/twirapp/twir/apps/bots/internal/discord/sended_messages_store"
	kickchat "github.com/twirapp/twir/apps/bots/internal/kick"
	"github.com/twirapp/twir/apps/bots/internal/messagehandler"
	mod_task_queue "github.com/twirapp/twir/apps/bots/internal/mod-task-queue"
	"github.com/twirapp/twir/apps/bots/internal/moderationhelpers"
	botplatforms "github.com/twirapp/twir/apps/bots/internal/platforms"
	"github.com/twirapp/twir/apps/bots/internal/services/channel"
	chattranslationsservice "github.com/twirapp/twir/apps/bots/internal/services/chat_translations"
	"github.com/twirapp/twir/apps/bots/internal/services/giveaways"
	"github.com/twirapp/twir/apps/bots/internal/services/keywords"
	toxicity_check "github.com/twirapp/twir/apps/bots/internal/services/toxicity-check"
	"github.com/twirapp/twir/apps/bots/internal/services/tts"
	"github.com/twirapp/twir/apps/bots/internal/services/voteban"
	"github.com/twirapp/twir/apps/bots/internal/services/ytsr"
	stream_handlers "github.com/twirapp/twir/apps/bots/internal/stream-handlers"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	vkchat "github.com/twirapp/twir/apps/bots/internal/vk"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	youtubechat "github.com/twirapp/twir/apps/bots/internal/youtube"
	"github.com/twirapp/twir/apps/bots/pkg/tlds"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	channelsgamesvotebancache "github.com/twirapp/twir/libs/cache/channels_games_voteban"
	channelsmoderationsettingscache "github.com/twirapp/twir/libs/cache/channels_moderation_settings"
	chattranslationssettingscache "github.com/twirapp/twir/libs/cache/chat_translations_settings"
	chatwallcacher "github.com/twirapp/twir/libs/cache/chat_wall"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	giveawayscache "github.com/twirapp/twir/libs/cache/giveaways"
	greetingscache "github.com/twirapp/twir/libs/cache/greetings"
	keywordscache "github.com/twirapp/twir/libs/cache/keywords"
	rolescache "github.com/twirapp/twir/libs/cache/roles"
	ttscache "github.com/twirapp/twir/libs/cache/tts"
	"github.com/twirapp/twir/libs/cache/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsemotesusagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_emotes_usages/datasources/clickhouse"
	channelsgamesvotebanrepository "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	channelsgamesvotebanpgx "github.com/twirapp/twir/libs/repositories/channels_games_voteban/pgx"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	channelsintegrationsdiscordpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_discord/datasource/postgres"
	channelsmoderationsettingsrepository "github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	channelsmoderationsettingsrepositorypostgres "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/datasource/postgres"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/chat_messages/datasources/clickhouse"
	channelschattrenslationsrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	channelschattrenslationsrepositorypostgres "github.com/twirapp/twir/libs/repositories/chat_translation/datasource/postgres"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallrepositorypostgres "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"
	discordsendednotifications "github.com/twirapp/twir/libs/repositories/discord_sended_notifications"
	discordsendednotificationspgx "github.com/twirapp/twir/libs/repositories/discord_sended_notifications/pgx"
	giveawaysrepository "github.com/twirapp/twir/libs/repositories/giveaways"
	giveawaysrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways/pgx"
	giveawaysparticipantsrepository "github.com/twirapp/twir/libs/repositories/giveaways_participants"
	giveawaysparticipantsrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways_participants/pgx"
	greetingsrepository "github.com/twirapp/twir/libs/repositories/greetings"
	greetingsrepositorypgx "github.com/twirapp/twir/libs/repositories/greetings/pgx"
	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
	keywordsrepositorypgx "github.com/twirapp/twir/libs/repositories/keywords/pgx"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	kickbotsrepositorypgx "github.com/twirapp/twir/libs/repositories/kick_bots/pgx"
	notificationsrepository "github.com/twirapp/twir/libs/repositories/notifications"
	notificationsrepositorypostgres "github.com/twirapp/twir/libs/repositories/notifications/datasource/postgres"
	overlays_tts_repository "github.com/twirapp/twir/libs/repositories/overlays_tts"
	overlays_tts_pgx "github.com/twirapp/twir/libs/repositories/overlays_tts/pgx"
	quotesrepository "github.com/twirapp/twir/libs/repositories/quotes"
	quotesrepositorypgx "github.com/twirapp/twir/libs/repositories/quotes/pgx"
	rolesrepository "github.com/twirapp/twir/libs/repositories/roles"
	rolesrepositorypgx "github.com/twirapp/twir/libs/repositories/roles/pgx"
	sentmessagesrepository "github.com/twirapp/twir/libs/repositories/sentmessages"
	sentmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/sentmessages/pgx"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	toxicmessagesrepository "github.com/twirapp/twir/libs/repositories/toxic_messages"
	toxicmessagesrepositorypgx "github.com/twirapp/twir/libs/repositories/toxic_messages/pgx"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	usersstatsrepository "github.com/twirapp/twir/libs/repositories/users_stats"
	usersstatsrepositorypostgres "github.com/twirapp/twir/libs/repositories/users_stats/datasources/postgres"
	vkvideobotsrepository "github.com/twirapp/twir/libs/repositories/vk_video_bots"
	vkvideobotsrepositorypgx "github.com/twirapp/twir/libs/repositories/vk_video_bots/datasource/postgres"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/wsrouter"
)

const Service = "bots"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,

	keywordsrepositorypgx.NewFx,
	wire.Bind(new(keywordsrepository.Repository), new(*keywordsrepositorypgx.Pgx)),
	quotesrepositorypgx.NewFx,
	wire.Bind(new(quotesrepository.Repository), new(*quotesrepositorypgx.Pgx)),
	greetingsrepositorypgx.NewFx,
	wire.Bind(new(greetingsrepository.Repository), new(*greetingsrepositorypgx.Pgx)),
	sentmessagesrepositorypgx.NewFx,
	wire.Bind(new(sentmessagesrepository.Repository), new(*sentmessagesrepositorypgx.Pgx)),
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	streamsrepositorypostgres.NewFx,
	wire.Bind(new(streamsrepository.Repository), new(*streamsrepositorypostgres.Pgx)),
	toxicmessagesrepositorypgx.NewFx,
	wire.Bind(new(toxicmessagesrepository.Repository), new(*toxicmessagesrepositorypgx.Pgx)),
	chatmessagesrepositoryclickhouse.NewFx,
	wire.Bind(new(chatmessagesrepository.Repository), new(*chatmessagesrepositoryclickhouse.Clickhouse)),
	channelscommandsprefixpgx.NewFx,
	wire.Bind(new(channelscommandsprefixrepository.Repository), new(*channelscommandsprefixpgx.Pgx)),
	chatwallrepositorypostgres.NewFx,
	wire.Bind(new(chatwallrepository.Repository), new(*chatwallrepositorypostgres.Pgx)),
	giveawaysrepositorypgx.NewFx,
	wire.Bind(new(giveawaysrepository.Repository), new(*giveawaysrepositorypgx.Pgx)),
	channelsmoderationsettingsrepositorypostgres.NewFx,
	wire.Bind(new(channelsmoderationsettingsrepository.Repository), new(*channelsmoderationsettingsrepositorypostgres.Pgx)),
	channelsemotesusagesrepositoryclickhouse.NewFx,
	wire.Bind(new(channelsemotesusagesrepository.Repository), new(*channelsemotesusagesrepositoryclickhouse.Clickhouse)),
	usersrepositorypgx.NewFx,
	wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
	usersstatsrepositorypostgres.NewFx,
	wire.Bind(new(usersstatsrepository.Repository), new(*usersstatsrepositorypostgres.Pgx)),
	rolesrepositorypgx.NewFx,
	wire.Bind(new(rolesrepository.Repository), new(*rolesrepositorypgx.Pgx)),
	overlays_tts_pgx.NewFx,
	wire.Bind(new(overlays_tts_repository.Repository), new(*overlays_tts_pgx.Pgx)),
	channelsgamesvotebanpgx.NewFx,
	wire.Bind(new(channelsgamesvotebanrepository.Repository), new(*channelsgamesvotebanpgx.Pgx)),
	channelschattrenslationsrepositorypostgres.NewFx,
	wire.Bind(new(channelschattrenslationsrepository.Repository), new(*channelschattrenslationsrepositorypostgres.Pgx)),
	giveawaysparticipantsrepositorypgx.NewFx,
	wire.Bind(new(giveawaysparticipantsrepository.Repository), new(*giveawaysparticipantsrepositorypgx.Pgx)),
	channelsintegrationsdiscordpostgres.NewFx,
	wire.Bind(new(channelsintegrationsdiscord.Repository), new(*channelsintegrationsdiscordpostgres.Pgx)),
	discordsendednotificationspgx.NewFx,
	wire.Bind(new(discordsendednotifications.Repository), new(*discordsendednotificationspgx.Pgx)),
	kickbotsrepositorypgx.NewFx,
	wire.Bind(new(kickbotsrepository.Repository), new(*kickbotsrepositorypgx.Pgx)),
	vkvideobotsrepositorypgx.NewFx,
	wire.Bind(new(vkvideobotsrepository.Repository), new(*vkvideobotsrepositorypgx.Pgx)),
	notificationsrepositorypostgres.NewFx,
	wire.Bind(new(notificationsrepository.Repository), new(*notificationsrepositorypostgres.Pgx)),

	tlds.New,
	channelservice.NewChannelService,
	NewWebsocketClient,
	workers.New,
	wire.Struct(new(workers.Opts), "*"),
	chatwallcacher.NewEnabledOnly,
	chatwallcacher.NewSettings,
	giveawayscache.New,
	mod_task_queue.NewRedisModTaskDistributor,
	wire.Bind(new(mod_task_queue.TaskDistributor), new(*mod_task_queue.ModTaskDistributor)),
	rolescache.New,
	toxicity_check.New,
	wire.Struct(new(toxicity_check.Opts), "*"),
	NewCommandsPrefixCache,
	ttscache.NewTTSSettings,
	keywordscache.New,
	greetingscache.New,
	channelcache.NewByTwitchUserID,
	twitchactions.New,
	wire.Struct(new(twitchactions.Opts), "*"),
	kickchat.NewChatClient,
	youtubechat.NewChatClient,
	NewVKVideoChatClient,
	vkchat.NewChatClient,
	botplatforms.NewChatRegistry,
	channelsmoderationsettingscache.New,
	channelsgamesvotebancache.New,
	moderationhelpers.New,
	wire.Struct(new(moderationhelpers.Opts), "*"),
	messagehandler.New,
	wire.Struct(new(messagehandler.Opts), "*"),
	channel.New,
	wire.Struct(new(channel.Opts), "*"),
	keywords.New,
	wire.Struct(new(keywords.Opts), "*"),
	tts.New,
	wire.Struct(new(tts.Opts), "*"),
	voteban.New,
	wire.Struct(new(voteban.Opts), "*"),
	chattranslationssettingscache.New,
	chattranslationsservice.New,
	wire.Struct(new(chattranslationsservice.Opts), "*"),
	giveaways.New,
	wire.Struct(new(giveaways.Opts), "*"),
	twitch.New,
	sended_messages_store.New,
	wire.Struct(new(sended_messages_store.Opts), "*"),
	discordmessagesupdater.New,
	wire.Struct(new(discordmessagesupdater.Opts), "*"),
	discord_go.New,
	wire.Struct(new(discord_go.Opts), "*"),
	wsrouter.NewNatsWsRouter,
	wire.Bind(new(wsrouter.WsRouter), new(*wsrouter.WsRouterNats)),

	wire.Struct(new(ytsr.Opts), "*"),
	RegisterYTSR,
	wire.Struct(new(mod_task_queue.RedisTaskProcessorOpts), "*"),
	mod_task_queue.NewRedisTaskProcessor,
	StartMetrics,
	wire.Struct(new(stream_handlers.Opts), "*"),
	RegisterStreamHandlers,
	wire.Struct(new(bus_listener.Opts), "*"),
	bus_listener.New,
	wire.Struct(new(discordbushandler.Opts), "*"),
	RegisterDiscordBusHandler,
	wire.Struct(new(notificationssync.Opts), "*"),
	notificationssync.New,
	wire.Struct(new(ApplicationDeps), "*"),
	NewApplication,
)

type YTSRRegistration struct{}
type MetricsServer struct{}
type StreamHandlersRegistration struct{}
type DiscordBusHandlerRegistration struct{}

type ApplicationDeps struct {
	Lifecycle         *lifecycle.Lifecycle
	YTSR              YTSRRegistration
	TaskProcessor     *mod_task_queue.RedisTaskProcessor
	Metrics           MetricsServer
	StreamHandlers    StreamHandlersRegistration
	BusListener       *bus_listener.BusListener
	DiscordBusHandler DiscordBusHandlerRegistration
	NotificationsSync *notificationssync.Service
	Logger            *slog.Logger
}

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewWebsocketClient(config cfg.Config) websockets.WebsocketClient {
	return clients.NewWebsocket(config.AppEnv)
}

func NewCommandsPrefixCache(
	repository channelscommandsprefixrepository.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix] {
	return channelscommandsprefixcache.New(repository, bus)
}

func RegisterYTSR(opts ytsr.Opts) (YTSRRegistration, error) {
	if err := ytsr.New(opts); err != nil {
		return YTSRRegistration{}, fmt.Errorf("register YTSR: %w", err)
	}

	return YTSRRegistration{}, nil
}

func StartMetrics(config cfg.Config) MetricsServer {
	if config.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go func() { _ = http.ListenAndServe("0.0.0.0:3000", nil) }()
	}

	return MetricsServer{}
}

func RegisterStreamHandlers(opts stream_handlers.Opts) StreamHandlersRegistration {
	stream_handlers.New(opts)
	return StreamHandlersRegistration{}
}

func RegisterDiscordBusHandler(opts discordbushandler.Opts) (DiscordBusHandlerRegistration, error) {
	if err := discordbushandler.New(opts); err != nil {
		return DiscordBusHandlerRegistration{}, fmt.Errorf("register Discord bus handler: %w", err)
	}

	return DiscordBusHandlerRegistration{}, nil
}

func NewApplication(deps ApplicationDeps) *Application {
	deps.Logger.Info("🚀 Bots started")
	return &Application{lifecycle: deps.Lifecycle}
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
