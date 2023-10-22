package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/events/internal"
	"github.com/satont/twir/apps/events/internal/grpc_impl"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/events"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
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

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(5)
	d.SetConnMaxIdleTime(1 * time.Minute)

	redisParams, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisParams)

	services := &internal.Services{
		DB:             db,
		Logger:         logger,
		Cfg:            cfg,
		BotsGrpc:       clients.NewBots(cfg.AppEnv),
		TokensGrpc:     clients.NewTokens(cfg.AppEnv),
		WebsocketsGrpc: clients.NewWebsocket(cfg.AppEnv),
		Redis:          redisClient,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.EVENTS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	events.RegisterEventsServer(grpcServer, grpc_impl.NewEvents(services))
	go grpcServer.Serve(lis)

	logger.Sugar().Info("Events microservices started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
	grpcServer.Stop()
}
