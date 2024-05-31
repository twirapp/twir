package baseapp

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func CreateBaseApp(appName string) fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewFx,
			twirsentry.NewFx(twirsentry.NewFxOpts{Service: appName}),
			logger.NewFx(
				logger.Opts{
					Service: appName,
				},
			),
			newRedis,
			newGorm,
			uptrace.NewFx(appName),
		),
	)
}

func newRedis(cfg config.Config) (*redis.Client, error) {
	redisOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(redisOpts)
	return redisClient, nil
}

func newGorm(cfg config.Config, lc fx.Lifecycle) (*gorm.DB, error) {
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             100 * time.Millisecond,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl),
		&gorm.Config{
			Logger: newLogger,
		},
	)
	if err != nil {
		return nil, err
	}
	d, _ := db.DB()
	d.SetMaxIdleConns(1)
	d.SetMaxOpenConns(10)
	d.SetConnMaxLifetime(time.Hour)

	lc.Append(
		fx.Hook{
			OnStop: func(_ context.Context) error {
				return d.Close()
			},
		},
	)

	return db, nil
}
