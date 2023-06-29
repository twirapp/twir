package main

import (
	"context"
	"github.com/satont/twir/apps/scheduler/internal/timers"
	"github.com/satont/twir/apps/scheduler/internal/types"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
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

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
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

	logger.Sugar().Info("Scheduler started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	cancelCtx()
	logger.Sugar().Info("Closing...")
}
