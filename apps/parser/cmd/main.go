package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/getsentry/sentry-go"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	commands_bus "github.com/satont/twir/apps/parser/internal/commands-bus"
	"github.com/satont/twir/apps/parser/internal/nats"
	chatwallservice "github.com/satont/twir/apps/parser/internal/services/chat_wall"
	task_queue "github.com/satont/twir/apps/parser/internal/task-queue"
	variables_bus "github.com/satont/twir/apps/parser/internal/variables-bus"
	cfg "github.com/satont/twir/libs/config"
	buscore "github.com/twirapp/twir/libs/bus-core"
	seventv "github.com/twirapp/twir/libs/cache/7tv"
	channelscommandsprefixcache "github.com/twirapp/twir/libs/cache/channels_commands_prefix"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	ttscache "github.com/twirapp/twir/libs/cache/tts"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/parser"
	channelscategoriesaliasespgx "github.com/twirapp/twir/libs/repositories/channels_categories_aliases/datasource/postgres"
	channelscommandsprefixpgx "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/pgx"
	scheduledvipsrepositorypgx "github.com/twirapp/twir/libs/repositories/scheduled_vips/datasource/postgres"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"

	chatwallcache "github.com/twirapp/twir/libs/cache/chat_wall"
	chatwallpgx "github.com/twirapp/twir/libs/repositories/chat_wall/datasource/postgres"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/parser/internal/commands"
	"github.com/satont/twir/apps/parser/internal/grpc_impl"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
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

	// gorm
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl))
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

	pgxconn, err := pgxpool.New(context.Background(), config.DatabaseUrl)
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

	tokensGrpc := clients.NewTokens(config.AppEnv)

	taskQueueDistributor := task_queue.NewRedisTaskDistributor(config, logger)
	queueProcessor := task_queue.NewRedisTaskProcessor(
		task_queue.RedisTaskProcessorOpts{
			Cfg:        *config,
			Logger:     logger,
			Gorm:       db,
			TokensGrpc: tokensGrpc,
		},
	)
	defer queueProcessor.Stop()

	go func() {
		err := queueProcessor.Start()
		if err != nil {
			logger.Fatal("Error starting queue processor", zap.Error(err))
		}
	}()

	bus := buscore.NewNatsBus(nc)

	redSync := redsync.New(goredis.NewPool(redisClient))

	trmManager, err := manager.New(trmpgx.NewDefaultFactory(pgxconn))
	if err != nil {
		panic(err)
	}

	commandsPrefixRepo := channelscommandsprefixpgx.New(channelscommandsprefixpgx.Opts{PgxPool: pgxconn})
	commandsPrefixRepoCache := channelscommandsprefixcache.New(commandsPrefixRepo, redisClient)
	ttsSettingsCacher := ttscache.NewTTSSettings(db, redisClient)
	spotifyRepo := channelsintegrationsspotifypgx.New(channelsintegrationsspotifypgx.Opts{PgxPool: pgxconn})
	usersRepo := usersrepositorypgx.New(usersrepositorypgx.Opts{PgxPool: pgxconn})
	channelsCategoriesAliasesRepo := channelscategoriesaliasespgx.New(channelscategoriesaliasespgx.Opts{PgxPool: pgxconn})
	scheduledVipsRepo := scheduledvipsrepositorypgx.New(scheduledvipsrepositorypgx.Opts{PgxPool: pgxconn})
	chatWallRepository := chatwallpgx.New(chatwallpgx.Opts{PgxPool: pgxconn})

	cachedTwitchClient, err := twitch.New(*config, tokensGrpc, redisClient)
	if err != nil {
		panic(err)
	}

	chatWallCache := chatwallcache.NewEnabledOnly(chatWallRepository, redisClient)

	chatWallService := chatwallservice.New(
		chatwallservice.Opts{
			ChatWallRepository: chatWallRepository,
			Gorm:               db,
			ChatWallCache:      chatWallCache,
			Redis:              redisClient,
			Config:             *config,
			TokensClient:       tokensGrpc,
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
			Dota:       clients.NewDota(config.AppEnv),
			Tokens:     tokensGrpc,
			Events:     clients.NewEvents(config.AppEnv),
			Ytsr:       clients.NewYtsr(config.AppEnv),
		},
		TaskDistributor:          taskQueueDistributor,
		Bus:                      bus,
		CommandsCache:            commandscache.New(db, redisClient),
		ChatWallRepo:             chatWallRepository,
		ChatWallCache:            chatWallCache,
		ChatWallService:          chatWallService,
		SevenTvCache:             seventv.New(redisClient),
		SevenTvCacheBySevenTvID:  seventv.NewBySeventvID(redisClient),
		RedSync:                  redSync,
		CommandsPrefixCache:      commandsPrefixRepoCache,
		CommandsPrefixRepository: commandsPrefixRepo,
		TTSCache:                 ttsSettingsCacher,
		SpotifyRepo:              spotifyRepo,
		UsersRepo:                usersRepo,
		CategoriesAliasesRepo:    channelsCategoriesAliasesRepo,
		ScheduledVipsRepo:        scheduledVipsRepo,
		CacheTwitchClient:        cachedTwitchClient,
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

	cmdBus := commands_bus.New(bus, s, commandsService, variablesService)
	cmdBus.Subscribe()
	defer cmdBus.Unsubscribe()

	variablesBus := variables_bus.New(bus, variablesService)
	variablesBus.Subscribe()
	defer variablesBus.Unsubscribe()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.PARSER_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	defer grpcServer.GracefulStop()
	parser.RegisterParserServer(
		grpcServer,
		grpc_impl.NewServer(s, commandsService, variablesService),
	)
	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	logger.Info("Parser microservice started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Sugar().Info("Exiting")
	appCtxCancel()
}
