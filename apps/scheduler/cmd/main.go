package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satont/twir/apps/scheduler/internal/grpc_impl"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/scheduler"
	"google.golang.org/grpc"

	s "github.com/satont/twir/apps/scheduler/internal/services"
	"github.com/satont/twir/apps/scheduler/internal/timers"
	"github.com/satont/twir/apps/scheduler/internal/types"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
	d, _ := db.DB()
	d.SetMaxOpenConns(5)
	d.SetConnMaxIdleTime(1 * time.Minute)
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}

	appCtx, cancelCtx := context.WithCancel(context.Background())

	pb, err := pubsub.NewPubSub(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}

	parserGrpc := clients.NewParser(cfg.AppEnv)
	commands := s.NewCommands(db, parserGrpc)
	roles := s.NewRoles(db)

	services := &types.Services{
		Grpc: &types.GrpcServices{
			Emotes: clients.NewEmotesCacher(cfg.AppEnv),
			Parser: parserGrpc,
			Tokens: clients.NewTokens(cfg.AppEnv),
		},
		Gorm:     db,
		Config:   cfg,
		PubSub:   pb,
		Commands: commands,
		Roles:    roles,
	}

	timers.NewWatched(appCtx, services)
	timers.NewEmotes(appCtx, services)
	timers.NewOnlineUsers(appCtx, services)
	timers.NewStreams(appCtx, services)
	timers.NewCommandsAndRoles(appCtx, services)
	timers.NewBannerChannels(appCtx, services)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.SCHEDULER_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	scheduler.RegisterSchedulerServer(
		grpcServer,
		&grpc_impl.SchedulerGrpc{
			Commands: commands,
			Roles:    roles,
		},
	)
	go grpcServer.Serve(lis)

	logger.Sugar().Info("Scheduler started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	cancelCtx()
	logger.Sugar().Info("Closing...")
}
