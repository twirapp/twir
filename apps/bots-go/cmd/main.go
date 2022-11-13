package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfg "github.com/satont/tsuwari/libs/config"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/getsentry/sentry-go"
	helix "github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/bots/internal/bots"
	nats_handlers "github.com/satont/tsuwari/apps/bots/internal/nats"
	myNats "github.com/satont/tsuwari/libs/nats"
	botsProto "github.com/satont/tsuwari/libs/nats/bots"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	logger, _ := zap.NewDevelopment()
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		logger.Sugar().Error(err)
		panic("Cannot load config of application")
	}

	if cfg.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDsn,
			Environment:      cfg.AppEnv,
			Debug:            true,
			TracesSampleRate: 1.0,
		})
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	natsEncodedConn, natsConn, err := myNats.New(cfg.NatsUrl)
	if err != nil {
		panic(err)
	}

	twitchClient := twitch.NewClient(&helix.Options{
		ClientID:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
		RedirectURI:  cfg.TwitchCallbackUrl,
	})

	botsService := bots.NewBotsService(&bots.NewBotsOpts{
		Twitch: twitchClient,
		DB:     db,
		Logger: logger,
		Cfg:    cfg,
		Nats:   natsConn,
	})
	natsHandlers := nats_handlers.NewNatsHandlers(&nats_handlers.NatsHandlersOpts{
		Db:          db,
		BotsService: botsService,
		Logger:      logger,
	})

	natsConn.Subscribe("bots.deleteMessages", natsHandlers.DeleteMessages)
	natsConn.Subscribe("bots.sendMessage", natsHandlers.SendMessage)
	natsConn.Subscribe(botsProto.SUBJECTS_JOIN_OR_LEAVE, natsHandlers.JoinOrLeave)
	natsConn.Subscribe("user.update", natsHandlers.UserUpdate)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
	natsEncodedConn.Close()
	natsConn.Close()
	d, _ = db.DB()
	d.Close()
}
