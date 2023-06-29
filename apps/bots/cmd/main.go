package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"github.com/satont/twir/apps/bots/internal/di"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/servers"
	"github.com/satont/twir/libs/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/getsentry/sentry-go"
	"github.com/satont/twir/apps/bots/internal/bots"
	"github.com/satont/twir/apps/bots/internal/grpc_impl"
	"github.com/satont/twir/apps/bots/internal/handlers"
	botsgrpc "github.com/satont/twir/libs/grpc/generated/bots"
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

	if cfg.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	if cfg.SentryDsn != "" {
		sentry.Init(
			sentry.ClientOptions{
				Dsn:              cfg.SentryDsn,
				Environment:      cfg.AppEnv,
				Debug:            true,
				TracesSampleRate: 1.0,
			},
		)
	}

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
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

	do.ProvideValue[tokens.TokensClient](di.Provider, clients.NewTokens(cfg.AppEnv))
	do.ProvideValue[events.EventsClient](di.Provider, clients.NewEvents(cfg.AppEnv))

	redisUrl, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisUrl)
	do.ProvideValue[redis.Client](di.Provider, *redisClient)

	botsService := bots.NewBotsService(
		&bots.NewBotsOpts{
			DB:         db,
			Logger:     logger,
			Cfg:        cfg,
			ParserGrpc: clients.NewParser(cfg.AppEnv),
		},
	)

	grpcNetListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.BOTS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge: 1 * time.Minute,
			},
		),
	)
	botsgrpc.RegisterBotsServer(
		grpcServer, grpc_impl.NewServer(
			&grpc_impl.GrpcImplOpts{
				Db:          db,
				BotsService: botsService,
				Logger:      logger,
				Cfg:         cfg,
			},
		),
	)
	go grpcServer.Serve(grpcNetListener)

	pb.Subscribe(
		"user.update", func(data []byte) {
			handlers.UserUpdate(db, botsService, data)
		},
	)
	pb.Subscribe(
		"stream.online", func(data []byte) {
			handlers.StreamsOnline(db, data)
		},
	)
	pb.Subscribe(
		"stream.offline", func(data []byte) {
			handlers.StreamsOffline(db, data)
		},
	)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
	grpcServer.Stop()
}
