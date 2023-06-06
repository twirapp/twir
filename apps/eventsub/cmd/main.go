package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/eventsub/internal/client"
	"github.com/satont/tsuwari/apps/eventsub/internal/grpm_impl"
	"github.com/satont/tsuwari/apps/eventsub/internal/handler"
	"github.com/satont/tsuwari/apps/eventsub/internal/helpers"
	"github.com/satont/tsuwari/apps/eventsub/internal/types"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"github.com/satont/tsuwari/libs/pubsub"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	appCtx, appCtxCancel := context.WithCancel(context.Background())

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	cfg, err := config.New()
	if err != nil {
		panic(err)
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

	appTun, err := helpers.GetAppTunnel(appCtx, cfg)
	if err != nil {
		panic(err)
	}

	appAddr := lo.
		If(cfg.AppEnv != "production", appTun.Addr().String()).
		Else(fmt.Sprintf("eventsub.%s", cfg.HostName))

	pb, err := pubsub.NewPubSub(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}

	redisUrl, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(redisUrl)

	services := &types.Services{
		Gorm:   db,
		Config: cfg,
		Grpc: &types.GrpcClients{
			Tokens: clients.NewTokens(cfg.AppEnv),
			Events: clients.NewEvents(cfg.AppEnv),
			Bots:   clients.NewBots(cfg.AppEnv),
			Parser: clients.NewParser(cfg.AppEnv),
		},
		PubSub: pb,
		Redis:  redisClient,
	}

	eventSubHandler := handler.NewHandler(services)
	go func() {
		if err := http.Serve(appTun, eventSubHandler.Manager); err != nil {
			panic(err)
		}
	}()

	eventSubClient, err := client.NewClient(appCtx, services, fmt.Sprintf("https://%s", appAddr))
	if err != nil {
		panic(err)
	}

	grpcImpl := grpm_impl.NewGrpcImpl(eventSubClient, services, fmt.Sprintf("https://%s", appAddr))
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.EVENTSUB_SERVER_PORT))
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionAge: 1 * time.Minute,
	}))
	eventsub.RegisterEventSubServer(grpcServer, grpcImpl)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	logger.Sugar().Infow("EventSub started", "addr", appAddr)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	appCtxCancel()
	appTun.Close()
	d.Close()
	logger.Sugar().Info("Closing...")
}
