package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/exaring/otelpgx"
	"github.com/getsentry/sentry-go"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	kvredis "github.com/twirapp/kv/stores/redis"
	commands_bus "github.com/twirapp/twir/apps/parser/internal/commands-bus"
	"github.com/twirapp/twir/apps/parser/internal/nats"
	chatwallservice "github.com/twirapp/twir/apps/parser/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/parser/internal/services/shortenedurls"
	variables_bus "github.com/twirapp/twir/apps/parser/internal/variables-bus"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/executron"
	"github.com/twirapp/twir/libs/baseapp"
	buscore "github.com/twirapp/twir/libs/bus-core"
	seventv "github.com/twirapp/twir/libs/cache/7tv"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	ttscache "github.com/twirapp/twir/libs/cache/tts"
	"github.com/twirapp/twir/libs/cache/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/i18n"
	channelscategoriesaliasespgx "github.com/twirapp/twir/libs/repositories/channels_categories_aliases/datasource/postgres"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	channelscommandsusagesclickhouse "github.com/twirapp/twir/libs/repositories/channels_commands_usages/datasources/clickhouse"
	channelsemotesusagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/channels_emotes_usages/datasources/clickhouse"
	channelseventslistpostgres "github.com/twirapp/twir/libs/repositories/channels_events_list/datasources/postgres"
	channelsgamesvotebanpgx "github.com/twirapp/twir/libs/repositories/channels_games_voteban/pgx"
	channelsgamesvotebanprogressstateredis "github.com/twirapp/twir/libs/repositories/channels_games_voteban_progress_state/redis"
	channelsinfohistorypostgres "github.com/twirapp/twir/libs/repositories/channels_info_history/datasource/postgres"
	channelsintegrationslastfmpostgres "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm/datasources/postgres"
	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	chatmessagesrepositoryclickhouse "github.com/twirapp/twir/libs/repositories/chat_messages/datasources/clickhouse"
	commandswithgroupsandresponsespostgres "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	overlaysttspgx "github.com/twirapp/twir/libs/repositories/overlays_tts/pgx"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	streamsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	"github.com/twirapp/twir/libs/uptrace"

	shortenedurlspgx "github.com/twirapp/twir/libs/repositories/shortened_urls/datasource/postgres"

	chatwallcache "github.com/twirapp/twir/libs/cache/chat_wall"
	chatwallpgx "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/parser/internal/commands"
	"github.com/twirapp/twir/apps/parser/internal/services/tts"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/apps/parser/internal/variables"
	"go.uber.org/zap"
)

func main() {
	appCtx, appCtxCancel := context.WithCancel(context.Background())

	config, err := cfg.New()
	if err != nil || config == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	if config.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	translationService, err := i18n.New(
		i18n.Opts{
			Store:         locales.Store,
			DefaultLocale: "en",
		},
	)
	if err != nil {
		panic(err)
	}

	if config.SentryDsn != "" {
		sentry.Init(
			sentry.ClientOptions{
				Dsn:              config.SentryDsn,
				Environment:      config.AppEnv,
				Debug:            false,
				TracesSampleRate: 1.0,
			},
		)
	}

	uptrace.New(*config, "parser")

	var logger *zap.Logger

	if config.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	zap.ReplaceGlobals(logger)

	connConfig, err := pgxpool.ParseConfig(config.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	connConfig.ConnConfig.Tracer = otelpgx.NewTracer()
	connConfig.MaxConnLifetime = time.Hour
	connConfig.MaxConnIdleTime = 5 * time.Minute
	connConfig.MaxConns = 100
	connConfig.MinConns = 1
	connConfig.HealthCheckPeriod = time.Minute
	connConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pgxconn, err := pgxpool.NewWithConfig(
		context.Background(),
		connConfig,
	)

	sqlDb := stdlib.OpenDBFromPool(pgxconn)

	dialector := postgres.New(
		postgres.Config{
			Conn: sqlDb,
		},
	)

	// gorm
	db, err := gorm.Open(
		dialector, &gorm.Config{
			PrepareStmt:            false,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxIdleConns(1)
	d.SetMaxOpenConns(10)
	d.SetConnMaxLifetime(time.Hour)
	defer d.Close()

	nc, err := nats.New(nats.Opts{Config: *config})
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	// sqlx
	dbConnOpts, err := pq.ParseURL(config.DatabaseUrl)
	if err != nil {
		panic(fmt.Errorf("cannot parse postgres url connection: %w", err))
	}
	pgConn, err := sqlx.ConnectContext(appCtx, "postgres", dbConnOpts)
	defer pgConn.Close()
	if err != nil {
		log.Fatalln(err)
	}

	// redis
	url, err := redis.ParseURL(config.RedisUrl)

	if err != nil {
		panic("Wrong redis url")
	}

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     url.Addr,
			Password: url.Password,
			DB:       url.DB,
			Username: url.Username,
		},
	)
	defer redisClient.Close()

	redisClient.Conn()

	bus := buscore.NewNatsBus(nc)

	redSync := redsync.New(goredis.NewPool(redisClient))

	trmManager, err := manager.New(trmpgx.NewDefaultFactory(pgxconn))
	if err != nil {
		panic(err)
	}

	clickhouseCreator := baseapp.NewClickHouse("parser")
	clickhouseClient, err := clickhouseCreator(*config)
	if err != nil {
		panic(err)
	}

	kvStorageRedis := kvredis.New(redisClient)

	commandsPrefixRepo := channelscommandsprefixpgx.New(channelscommandsprefixpgx.Opts{PgxPool: pgxconn})
	commandsPrefixRepoCache := channelscommandsprefixcache.New(commandsPrefixRepo, bus)
	ttsRepository := overlaysttspgx.NewFx(pgxconn)
	ttsSettingsCacher := ttscache.NewTTSSettings(ttsRepository, kvStorageRedis)
	spotifyRepo := channelsintegrationsspotifypgx.New(channelsintegrationsspotifypgx.Opts{PgxPool: pgxconn})
	usersRepo := usersrepositorypgx.New(usersrepositorypgx.Opts{PgxPool: pgxconn})
	channelsCategoriesAliasesRepo := channelscategoriesaliasespgx.New(channelscategoriesaliasespgx.Opts{PgxPool: pgxconn})
	scheduledVipsRepo := scheduledvipsrepositorypgx.New(scheduledvipsrepositorypgx.Opts{PgxPool: pgxconn})
	chatWallRepository := chatwallpgx.New(chatwallpgx.Opts{PgxPool: pgxconn})
	channelsInfoHistoryRepo := channelsinfohistorypostgres.New(channelsinfohistorypostgres.Opts{PgxPool: pgxconn})
	shortenedUrlsRepo := shortenedurlspgx.New(shortenedurlspgx.Opts{PgxPool: pgxconn})
	streamsRepository := streamsrepositorypostgres.New(streamsrepositorypostgres.Opts{PgxPool: pgxconn})
	channelsEmotesUsage := channelsemotesusagesrepositoryclickhouse.New(channelsemotesusagesrepositoryclickhouse.Opts{Client: clickhouseClient})
	channelsCommandsUsagesRepo := channelscommandsusagesclickhouse.New(channelscommandsusagesclickhouse.Opts{Client: clickhouseClient})
	chatMessagesRepo := chatmessagesrepositoryclickhouse.New(chatmessagesrepositoryclickhouse.Opts{Client: clickhouseClient})
	channelsEventListRepo := channelseventslistpostgres.New(channelseventslistpostgres.Opts{PgxPool: pgxconn})
	channelsGamesVotebanRepo := channelsgamesvotebanpgx.New(channelsgamesvotebanpgx.Opts{PgxPool: pgxconn})
	channelsGamesVotebanProgressStateRepo := channelsgamesvotebanprogressstateredis.New(channelsgamesvotebanprogressstateredis.Opts{Redis: redisClient})
	lastfmRepo := channelsintegrationslastfmpostgres.New(channelsintegrationslastfmpostgres.Opts{PgxPool: pgxconn})

	cachedTwitchClient, err := twitch.New(*config, bus, redisClient)
	if err != nil {
		panic(err)
	}

	chatWallCache := chatwallcache.NewEnabledOnly(chatWallRepository, bus)

	chatWallService := chatwallservice.New(
		chatwallservice.Opts{
			ChatWallRepository: chatWallRepository,
			ChatMessagesRepo:   chatMessagesRepo,
			Gorm:               db,
			ChatWallCache:      chatWallCache,
			Redis:              redisClient,
			Config:             *config,
			TwirBus:            bus,
		},
	)

	s := &services.Services{
		Config:     config,
		Logger:     logger,
		Gorm:       db,
		Sqlx:       pgConn,
		Redis:      redisClient,
		TrmManager: trmManager,
		GrpcClients: &services.Grpc{
			WebSockets: clients.NewWebsocket(config.AppEnv),
		},
		Bus: bus,
		CommandsCache: commandscache.New(
			commandswithgroupsandresponsespostgres.New(commandswithgroupsandresponsespostgres.Opts{PgxPool: pgxconn}),
			bus,
		),
		ChatWallRepo:             chatWallRepository,
		ChatWallCache:            chatWallCache,
		ChatWallService:          chatWallService,
		SevenTvCache:             seventv.New(kvStorageRedis, *config),
		RedSync:                  redSync,
		CommandsPrefixCache:      commandsPrefixRepoCache,
		CommandsPrefixRepository: commandsPrefixRepo,
		TTSCache:                 ttsSettingsCacher,
		TTSRepository:            ttsRepository,
		TTSService:               tts.New(ttsRepository, config),
		SpotifyRepo:              spotifyRepo,
		UsersRepo:                usersRepo,
		CategoriesAliasesRepo:    channelsCategoriesAliasesRepo,
		ScheduledVipsRepo:        scheduledVipsRepo,
		CacheTwitchClient:        cachedTwitchClient,
		ChannelsInfoHistoryRepo:  channelsInfoHistoryRepo,
		ChannelEventListsRepo:    channelsEventListRepo,
		ShortUrlServices: shortenedurls.New(
			shortenedurls.Opts{
				Repository: shortenedUrlsRepo,
				Config:     *config,
			},
		),
		ChannelEmotesUsagesRepo:           channelsEmotesUsage,
		ChannelsCommandsUsagesRepo:        channelsCommandsUsagesRepo,
		ChatMessagesRepo:                  chatMessagesRepo,
		ChannelsGamesVotebanRepo:          channelsGamesVotebanRepo,
		ChannelsGamesVotebanProgressState: channelsGamesVotebanProgressStateRepo,
		LastfmRepo:                        lastfmRepo,
		Executron:                         executron.New(*config, redisClient),
		I18n:                              translationService,
	}

	variablesService := variables.New(
		&variables.Opts{
			Services: s,
		},
	)
	commandsService := commands.New(
		&commands.Opts{
			Services:         s,
			VariablesService: variablesService,
		},
	)

	cmdBus := commands_bus.New(bus, s, commandsService, variablesService, streamsRepository)
	cmdBus.Subscribe()
	defer cmdBus.Unsubscribe()

	variablesBus := variables_bus.New(bus, variablesService)
	variablesBus.Subscribe()
	defer variablesBus.Unsubscribe()

	logger.Info("Parser microservice started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Sugar().Info("Exiting")
	appCtxCancel()
}
