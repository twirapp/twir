package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/goforj/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/parser/internal/commands"
	commandsbus "github.com/twirapp/twir/apps/parser/internal/commands-bus"
	chatwallservice "github.com/twirapp/twir/apps/parser/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/parser/internal/services/shortenedurls"
	"github.com/twirapp/twir/apps/parser/internal/services/tts"
	parserservices "github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/apps/parser/internal/variables"
	variablesbus "github.com/twirapp/twir/apps/parser/internal/variables-bus"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	seventv "github.com/twirapp/twir/libs/cache/7tv"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	chatwallcache "github.com/twirapp/twir/libs/cache/chat_wall"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	quotescache "github.com/twirapp/twir/libs/cache/quotes"
	ttscache "github.com/twirapp/twir/libs/cache/tts"
	"github.com/twirapp/twir/libs/cache/twitch"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/i18n"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelscategoriesaliasesrepository "github.com/twirapp/twir/libs/repositories/channels_categories_aliases"
	channelscategoriesaliasespostgres "github.com/twirapp/twir/libs/repositories/channels_categories_aliases/datasource/postgres"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelscommandsusagesrepository "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
	channelscommandsusagesclickhouse "github.com/twirapp/twir/libs/repositories/channels_commands_usages/datasources/clickhouse"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	channelsemotesusagesclickhouse "github.com/twirapp/twir/libs/repositories/channels_emotes_usages/datasources/clickhouse"
	channelseventslistrepository "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistpostgres "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	channelsgamesvotebanrepository "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	channelsgamesvotebanpgx "github.com/twirapp/twir/libs/repositories/channels_games_voteban/pgx"
	channelsinfohistoryrepository "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsinfohistorypostgres "github.com/twirapp/twir/libs/repositories/channels_info_history/datasource/postgres"
	channelsintegrationslastfmrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	channelsintegrationslastfmpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm/datasources/postgres"
	channelsintegrationsspotifyrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	chatmessagesrepository "github.com/twirapp/twir/libs/repositories/chat_messages"
	chatmessagesclickhouse "github.com/twirapp/twir/libs/repositories/chat_messages/datasources/clickhouse"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallpostgres "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandswithgroupsandresponsespgx "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	faceitintegrationrepository "github.com/twirapp/twir/libs/repositories/faceit_integration"
	faceitintegrationpostgres "github.com/twirapp/twir/libs/repositories/faceit_integration/datasource/postgres"
	overlaysttsrepository "github.com/twirapp/twir/libs/repositories/overlays_tts"
	overlaysttspgx "github.com/twirapp/twir/libs/repositories/overlays_tts/pgx"
	quotesrepository "github.com/twirapp/twir/libs/repositories/quotes"
	quotespgx "github.com/twirapp/twir/libs/repositories/quotes/pgx"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipspostgres "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	shortenedurlspostgres "github.com/twirapp/twir/libs/repositories/shortened_urls/datasource/postgres"
	songrequestssettingsrepository "github.com/twirapp/twir/libs/repositories/song_requests_settings"
	songrequestssettingspostgres "github.com/twirapp/twir/libs/repositories/song_requests_settings/datasource/postgres"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamspostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	userswithstatsrepository "github.com/twirapp/twir/libs/repositories/userswithstats"
	userswithstatspostgres "github.com/twirapp/twir/libs/repositories/userswithstats/datasource/postgres"
	vkintegrationrepository "github.com/twirapp/twir/libs/repositories/vk_integration"
	vkintegrationpostgres "github.com/twirapp/twir/libs/repositories/vk_integration/datasource/postgres"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/zap"
)

const Service = "parser"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	channelscategoriesaliasespostgres.NewFx,
	wire.Bind(new(channelscategoriesaliasesrepository.Repository), new(*channelscategoriesaliasespostgres.Pgx)),
	channelscommandsprefixpgx.NewFx,
	wire.Bind(new(channelscommandsprefixrepository.Repository), new(*channelscommandsprefixpgx.Pgx)),
	channelscommandsusagesclickhouse.NewFx,
	wire.Bind(new(channelscommandsusagesrepository.Repository), new(*channelscommandsusagesclickhouse.Clickhouse)),
	channelsemotesusagesclickhouse.NewFx,
	wire.Bind(new(channelsemotesusagesrepository.Repository), new(*channelsemotesusagesclickhouse.Clickhouse)),
	channelseventslistpostgres.NewFx,
	wire.Bind(new(channelseventslistrepository.Repository), new(*channelseventslistpostgres.Pgx)),
	channelsgamesvotebanpgx.NewFx,
	wire.Bind(new(channelsgamesvotebanrepository.Repository), new(*channelsgamesvotebanpgx.Pgx)),
	channelsinfohistorypostgres.NewFx,
	wire.Bind(new(channelsinfohistoryrepository.Repository), new(*channelsinfohistorypostgres.Pgx)),
	channelsintegrationslastfmpostgres.NewFx,
	wire.Bind(new(channelsintegrationslastfmrepository.Repository), new(*channelsintegrationslastfmpostgres.Pgx)),
	channelsintegrationsspotifypgx.NewFx,
	wire.Bind(new(channelsintegrationsspotifyrepository.Repository), new(*channelsintegrationsspotifypgx.Pgx)),
	chatmessagesclickhouse.NewFx,
	wire.Bind(new(chatmessagesrepository.Repository), new(*chatmessagesclickhouse.Clickhouse)),
	chatwallpostgres.NewFx,
	wire.Bind(new(chatwallrepository.Repository), new(*chatwallpostgres.Pgx)),
	commandswithgroupsandresponsespgx.NewFx,
	wire.Bind(new(commandswithgroupsandresponsesrepository.Repository), new(*commandswithgroupsandresponsespgx.Pgx)),
	faceitintegrationpostgres.NewFx,
	wire.Bind(new(faceitintegrationrepository.Repository), new(*faceitintegrationpostgres.Pgx)),
	overlaysttspgx.NewFx,
	wire.Bind(new(overlaysttsrepository.Repository), new(*overlaysttspgx.Pgx)),
	quotespgx.NewFx,
	wire.Bind(new(quotesrepository.Repository), new(*quotespgx.Pgx)),
	scheduledvipspostgres.NewFx,
	wire.Bind(new(scheduledvipsrepository.Repository), new(*scheduledvipspostgres.Pgx)),
	shortenedurlspostgres.NewFx,
	wire.Bind(new(shortenedurlsrepository.Repository), new(*shortenedurlspostgres.Pgx)),
	streamspostgres.NewFx,
	wire.Bind(new(streamsrepository.Repository), new(*streamspostgres.Pgx)),
	usersrepositorypgx.NewFx,
	wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
	userswithstatspostgres.NewFx,
	wire.Bind(new(userswithstatsrepository.Repository), new(*userswithstatspostgres.Pgx)),
	vkintegrationpostgres.NewFx,
	wire.Bind(new(vkintegrationrepository.Repository), new(*vkintegrationpostgres.Pgx)),
	songrequestssettingspostgres.NewFx,
	wire.Bind(new(songrequestssettingsrepository.Repository), new(*songrequestssettingspostgres.Pgx)),
	channelscommandsprefixcache.New,
	commandscache.New,
	quotescache.New,
	ttscache.NewTTSSettings,
	chatwallcache.NewEnabledOnly,
	seventv.New,
	twitch.New,
	channelservice.NewChannelService,
	chatwallservice.New,
	shortenedurls.New,
	tts.New,
	NewI18n,
	NewConfigPointer,
	NewZapLogger,
	NewSqlx,
	NewRedSync,
	NewGrpc,
	wire.Struct(new(parserservices.Services),
		"Config", "Logger", "Gorm", "PgxPool", "Sqlx", "Redis", "GrpcClients", "Bus", "TrmManager",
		"CommandsCache", "CommandsPrefixCache", "SevenTvCache", "ChatWallCache", "ChatWallService", "RedSync",
		"ChannelsRepo", "ChannelService", "CommandsPrefixRepository", "TTSCache", "TTSRepository", "TTSService",
		"SpotifyRepo", "UsersRepo", "CategoriesAliasesRepo", "ScheduledVipsRepo", "CacheTwitchClient", "ChatWallRepo",
		"ChannelsInfoHistoryRepo", "ChannelEmotesUsagesRepo", "ChannelsCommandsUsagesRepo", "ChatMessagesRepo",
		"ChannelEventListsRepo", "ChannelsGamesVotebanRepo", "ShortUrlServices", "LastfmRepo", "VKRepo", "FaceitRepo",
		"I18n", "UsersWithStatsRepository", "QuotesRepo", "QuotesCacher",
		"SongRequestsSettingsRepo",
	),
	variables.New,
	commands.New,
	commandsbus.New,
	variablesbus.New,
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewConfigPointer(applicationConfig config.Config) *config.Config {
	return &applicationConfig
}

func NewZapLogger(applicationConfig config.Config) *zap.Logger {
	var logger *zap.Logger
	if applicationConfig.AppEnv == "development" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(logger)
	return logger
}

func NewSqlx(pool *pgxpool.Pool, lc *lifecycle.Lifecycle) *sqlx.DB {
	db := sqlx.NewDb(stdlib.OpenDBFromPool(pool), "postgres")
	lc.Append(lifecycle.Hook{OnStop: func(context.Context) error { return db.Close() }})
	return db
}

func NewRedSync(redisClient *redis.Client) *redsync.Redsync {
	return redsync.New(goredis.NewPool(redisClient))
}

func NewGrpc(applicationConfig config.Config) *parserservices.Grpc {
	return &parserservices.Grpc{WebSockets: clients.NewWebsocket(applicationConfig.AppEnv)}
}

func NewI18n() (*i18n.I18n, error) {
	return i18n.New(i18n.Opts{Store: locales.Store, DefaultLocale: "en"})
}

func NewApplication(
	lc *lifecycle.Lifecycle,
	logger *slog.Logger,
	applicationConfig config.Config,
	commandsBus *commandsbus.CommandsBus,
	variablesBus *variablesbus.VariablesBus,
) (*Application, error) {
	if err := commandsBus.Subscribe(); err != nil {
		return nil, fmt.Errorf("subscribe commands bus: %w", err)
	}
	if err := variablesBus.Subscribe(); err != nil {
		commandsBus.Unsubscribe()
		return nil, fmt.Errorf("subscribe variables bus: %w", err)
	}

	lc.Append(lifecycle.Hook{
		OnStop: func(context.Context) error {
			variablesErr := variablesBus.Unsubscribe()
			commandsBus.Unsubscribe()
			return variablesErr
		},
	})

	if applicationConfig.AppEnv != "development" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		metricsServer := &http.Server{Addr: "0.0.0.0:3000", Handler: mux}
		lc.Append(lifecycle.Hook{
			OnStart: func(context.Context) error {
				go func() {
					if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						logger.Error("metrics server stopped", slog.Any("error", err))
					}
				}()
				return nil
			},
			OnStop: metricsServer.Shutdown,
		})
	}

	logger.Info("Parser microservice started")
	return &Application{lifecycle: lc}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
