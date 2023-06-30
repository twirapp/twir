package main

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/twir/apps/timers/internal/di"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satont/twir/apps/timers/internal/scheduler"
	"github.com/satont/twir/apps/timers/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/servers"

	cfg "github.com/satont/twir/libs/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/satont/twir/apps/timers/internal/grpc_impl"
	timersgrpc "github.com/satont/twir/libs/grpc/generated/timers"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	logger, _ := zap.NewDevelopment()

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}

	parserGrpcClient := clients.NewParser(cfg.AppEnv)
	botsGrpcClient := clients.NewBots(cfg.AppEnv)
	do.ProvideValue[tokens.TokensClient](di.Provider, clients.NewTokens(cfg.AppEnv))

	scheduler := scheduler.New(cfg, db, logger, parserGrpcClient, botsGrpcClient)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.TIMERS_SERVER_PORT))
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
	timersgrpc.RegisterTimersServer(
		grpcServer, grpc_impl.New(
			&grpc_impl.TimersGrpcServerOpts{
				Db:        db,
				Logger:    logger,
				Scheduler: scheduler,
			},
		),
	)
	go grpcServer.Serve(lis)

	timers := []*model.ChannelsTimers{}
	err = db.Model(&model.ChannelsTimers{}).
		Where("1 = 1").
		Update("lastTriggerMessageNumber", 0).
		Error
	if err != nil {
		logger.Sugar().Error(err)
	}
	err = db.Preload("Responses").Preload("Channel").Find(&timers).Error

	if err != nil {
		panic(err)
	} else {
		for _, timer := range timers {
			if timer.Channel != nil && (!timer.Channel.IsEnabled || timer.Channel.IsBanned) {
				continue
			}

			if timer.Enabled {
				scheduler.AddTimer(
					&types.Timer{
						Model:     timer,
						SendIndex: 0,
					},
				)
			}
		}
	}

	logger.Sugar().Info("Timers microservice started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
	grpcServer.Stop()
	log.Fatalf("Exiting")
}
