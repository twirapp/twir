package main

import (
	"context"
	"github.com/satont/tsuwari/apps/scheduler/internal/timers"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"
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

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}

	appCtx := context.Background()
	services := &types.Services{
		Grpc: &types.GrpcServices{
			Emotes:  clients.NewEmotesCacher(cfg.AppEnv),
			Parser:  clients.NewParser(cfg.AppEnv),
			Tokens:  clients.NewTokens(cfg.AppEnv),
			Watched: clients.NewWatched(cfg.AppEnv),
		},
		Gorm:   db,
		Config: cfg,
	}

	timers.NewWatched(appCtx, services)
	timers.NewEmotes(appCtx, services)
	timers.NewOnlineUsers(appCtx, services)

	logger.Sugar().Info("Started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	appCtx.Done()
	logger.Sugar().Info("Closing...")
}
