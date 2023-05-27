package main

import (
	"context"
	"fmt"
	"github.com/satont/tsuwari/apps/scheduler/grpc_impl"
	"github.com/satont/tsuwari/apps/scheduler/internal/timers"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"github.com/satont/tsuwari/libs/pubsub"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}

	appCtx, cancelCtx := context.WithCancel(context.Background())

	pb, err := pubsub.NewPubSub(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}

	services := &types.Services{
		Grpc: &types.GrpcServices{
			Emotes:  clients.NewEmotesCacher(cfg.AppEnv),
			Parser:  clients.NewParser(cfg.AppEnv),
			Tokens:  clients.NewTokens(cfg.AppEnv),
			Watched: clients.NewWatched(cfg.AppEnv),
		},
		Gorm:   db,
		Config: cfg,
		PubSub: pb,
	}

	timers.NewWatched(appCtx, services)
	timers.NewEmotes(appCtx, services)
	timers.NewOnlineUsers(appCtx, services)
	timers.NewStreams(appCtx, services)
	cmds := timers.NewDefaultCommands(services)
	cmds.Run(appCtx)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.SCHEDULER_SERVER_PORT))
	if err != nil {
		logger.Sugar().Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	scheduler.RegisterSchedulerServer(
		grpcServer,
		grpc_impl.NewGrpcImpl(cmds, services),
	)
	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	logger.Sugar().Info("Scheduler started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	cancelCtx()
	logger.Sugar().Info("Closing...")
}
