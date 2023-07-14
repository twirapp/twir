package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/pubsub"

	"github.com/satont/twir/apps/scheduler/internal/timers"
	"github.com/satont/twir/apps/scheduler/internal/types"
)

func run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pb, err := pubsub.NewPubSub(cfg.RedisUrl)
	if err != nil {
		return fmt.Errorf("new pubsub: %w", err)
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

	timers.NewWatched(ctx, services)
	timers.NewEmotes(ctx, services)
	timers.NewOnlineUsers(ctx, services)
	timers.NewStreams(ctx, services)

	logger.Sugar().Info("Scheduler started")
	<-ctx.Done()
	logger.Sugar().Info("Closing...")

	return ctx.Err()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}
