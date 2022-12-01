package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"github.com/satont/tsuwari/libs/pubsub"
	"google.golang.org/grpc"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/getsentry/sentry-go"
	helix "github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/bots/internal/bots"
	"github.com/satont/tsuwari/apps/bots/internal/grpc_impl"
	"github.com/satont/tsuwari/apps/bots/internal/handlers"
	botsgrpc "github.com/satont/tsuwari/libs/grpc/generated/bots"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	var logger *zap.Logger
	cfg, err := cfg.New()

	if err != nil || cfg == nil {
		panic("Cannot load config of application")
	}

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}
	zap.ReplaceGlobals(logger)

	if cfg.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDsn,
			Environment:      cfg.AppEnv,
			Debug:            true,
			TracesSampleRate: 1.0,
		})
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	pb, err := pubsub.NewPubSub(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}

	twitchClient := twitch.NewClient(&helix.Options{
		ClientID:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
		RedirectURI:  cfg.TwitchCallbackUrl,
	})

	botsService := bots.NewBotsService(&bots.NewBotsOpts{
		Twitch:     twitchClient,
		DB:         db,
		Logger:     logger,
		Cfg:        cfg,
		ParserGrpc: clients.NewParser(cfg.AppEnv),
	})

	grpcNetListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.BOTS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	botsgrpc.RegisterBotsServer(grpcServer, grpc_impl.NewServer(&grpc_impl.GrpcImplOpts{
		Db:          db,
		BotsService: botsService,
		Logger:      logger,
	}))
	go grpcServer.Serve(grpcNetListener)

	pb.Subscribe("user.update", func(data string) {
		handlers.UserUpdate(db, botsService, data)
	})
	pb.Subscribe("streams.online", func(data string) {
		handlers.StreamsOnline(db, data)
	})
	pb.Subscribe("streams.offline", func(data string) {
		handlers.StreamsOffline(db, data)
	})

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")

	d, _ = db.DB()
	d.Close()
	grpcServer.Stop()
}
