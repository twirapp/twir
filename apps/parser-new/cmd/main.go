package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	appCtx, appCtxCancel := context.WithCancel(context.Background())

	config, err := cfg.New()
	if err != nil || config == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	if config.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	if config.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn:              config.SentryDsn,
			Environment:      config.AppEnv,
			Debug:            true,
			TracesSampleRate: 1.0,
		})
	}

	var logger *zap.Logger

	if config.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	zap.ReplaceGlobals(logger)

	// gorm
	db, err := gorm.Open(postgres.Open(config.DatabaseUrl))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	// sqlx
	dbConnOpts, err := pq.ParseURL(config.DatabaseUrl)
	if err != nil {
		panic(fmt.Errorf("cannot parse postgres url connection: %w", err))
	}
	pgConn, err := sqlx.Connect("postgres", dbConnOpts)
	if err != nil {
		log.Fatalln(err)
	}

	// redis
	url, err := redis.ParseURL(config.RedisUrl)

	if err != nil {
		panic("Wrong redis url")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     url.Addr,
		Password: url.Password,
		DB:       url.DB,
		Username: url.Username,
	})

	redisClient.Conn()

	services := &services.Services{
		Config: config,
		Logger: logger,
		Gorm:   db,
		Sqlx:   pgConn,
		Redis:  redisClient,
		GrpcClients: &services.ServicesGrpc{
			WebSockets: clients.NewWebsocket(config.AppEnv),
			Bots:       clients.NewBots(config.AppEnv),
			Dota:       clients.NewDota(config.AppEnv),
			Eval:       clients.NewEval(config.AppEnv),
			Tokens:     clients.NewTokens(config.AppEnv),
			Events:     clients.NewEvents(config.AppEnv),
			Ytsr:       clients.NewYtsr(config.AppEnv),
		},
	}

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	logger.Sugar().Info("Exiting")
	pgConn.Close()
	redisClient.Close()
	d.Close()
	appCtxCancel()
}
