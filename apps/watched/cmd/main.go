package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/satont/tsuwari/apps/watched/internal/handlers"
	cfg "github.com/satont/tsuwari/libs/config"
	myNats "github.com/satont/tsuwari/libs/nats"
	"github.com/satont/tsuwari/libs/nats/watched"
	"go.uber.org/zap"
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

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	natsEncodedConn, _, err := myNats.New(cfg.NatsUrl)
	if err != nil {
		panic(err)
	}

	h := handlers.NewHandlers(handlers.HandlersOpts{
		DB:     db,
		Cfg:    cfg,
		Logger: logger,
	})

	natsEncodedConn.QueueSubscribe(watched.SUBJECTS_PROCESS_WATCHED_STREAMS, "watched", h.ProcessWatchedStreams)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Fatalf("Exiting")
}
