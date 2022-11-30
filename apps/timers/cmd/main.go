package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/satont/tsuwari/apps/timers/internal/scheduler"
	"github.com/satont/tsuwari/apps/timers/internal/types"
	"google.golang.org/grpc"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/grpc/servers"

	cfg "github.com/satont/tsuwari/libs/config"

	twitch "github.com/satont/tsuwari/libs/twitch"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/timers/internal/grpc_impl"
	timersgrpc "github.com/satont/tsuwari/libs/grpc/generated/timers"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	logger, _ := zap.NewDevelopment()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}

	t := twitch.NewClient(&helix.Options{
		ClientID:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
	})

	parserGrpcClient := clients.NewParser(cfg.AppEnv)
	botsGrpcClient := clients.NewBots(cfg.AppEnv)

	scheduler := scheduler.New(cfg, t, db, logger, parserGrpcClient, botsGrpcClient)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.TIMERS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	timersgrpc.RegisterTimersServer(grpcServer, grpc_impl.New(&grpc_impl.TimersGrpcServerOpts{
		Db:        db,
		Logger:    logger,
		Scheduler: scheduler,
	}))
	go grpcServer.Serve(lis)

	timers := []*model.ChannelsTimers{}
	err = db.Model(&model.ChannelsTimers{}).
		Where("1 = 1").
		Update("lastTriggerMessageNumber", 0).
		Error
	if err != nil {
		logger.Sugar().Error(err)
	}
	err = db.Preload("Responses").Find(&timers).Error

	if err != nil {
		panic(err)
	} else {
		for _, timer := range timers {
			if timer.Enabled {
				scheduler.AddTimer(&types.Timer{
					Model:     timer,
					SendIndex: 0,
				})
			}
		}
	}

	logger.Sugar().Info("Started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	grpcServer.Stop()
	log.Fatalf("Exiting")
}
