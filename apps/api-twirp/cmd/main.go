package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/twirp_handler"
	cfg "github.com/satont/tsuwari/libs/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"net/http"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	config, err := cfg.New()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisOpts)

	twirpPathPrefix, twirpHandler := twirp_handler.New(twirp_handler.Opts{
		Redis: redisClient,
		DB:    db,
	})

	mux := http.NewServeMux()
	mux.Handle(twirpPathPrefix, twirpHandler)
	http.ListenAndServe(":3002", mux)
}
