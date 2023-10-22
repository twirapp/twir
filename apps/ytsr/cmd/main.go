package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/satont/twir/apps/ytsr/internal/grpc_impl"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/ytsr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	config, err := cfg.New()
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger
	if config.AppEnv != "production" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.YTSR_SERVER_PORT))
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
	ytsr.RegisterYtsrServer(grpcServer, grpc_impl.NewYtsrServer(*config, logger))
	go grpcServer.Serve(lis)

	logger.Sugar().Info("YTSR microservice started")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Fatalf("Exiting")
}
