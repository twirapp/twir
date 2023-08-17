package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/satont/twir/libs/grpc/clients"

	"github.com/satont/twir/apps/watched/internal/grpc_impl"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/watched"
	"github.com/satont/twir/libs/grpc/servers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	cfg, err := cfg.New()
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Error),
		},
	)
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(5)
	d.SetConnMaxIdleTime(1 * time.Minute)

	tokensGrpc := clients.NewTokens(cfg.AppEnv)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.WATCHED_SERVER_PORT))
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
	watched.RegisterWatchedServer(
		grpcServer, grpc_impl.New(
			&grpc_impl.WatchedGrpcServerOpts{
				Db:         db,
				Cfg:        cfg,
				Logger:     logger,
				TokensGrpc: tokensGrpc,
			},
		),
	)
	go grpcServer.Serve(lis)

	log.Println("Watched microservice started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Fatalf("Exiting")
}
