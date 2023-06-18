package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/sessions"
	"github.com/satont/tsuwari/apps/api-twirp/internal/twirp_handlers"
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
		logger.Sugar().Panic(err)
	}

	db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Panic("failed to connect database", err)
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		logger.Sugar().Panic(err)
	}
	redisClient := redis.NewClient(redisOpts)

	sessionManager := sessions.New(redisClient)

	mux := http.NewServeMux()

	twirpOpts := twirp_handlers.Opts{
		Redis:          redisClient,
		DB:             db,
		SessionManager: sessionManager,
	}

	mux.Handle(twirp_handlers.NewProtected(twirpOpts))
	mux.Handle(twirp_handlers.NewUnProtected(twirpOpts))

	logger.Sugar().Panic(http.ListenAndServe(":3002", sessionManager.LoadAndSave(mux)))
}
