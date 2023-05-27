package main

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-new/internal/http/fiber"
	"github.com/satont/tsuwari/apps/api-new/internal/http/middlewares"
	"github.com/satont/tsuwari/apps/api-new/internal/http/routes"
	"github.com/satont/tsuwari/apps/api-new/internal/http/routes/v1_handlers"
	config "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

func main() {
	z, _ := zap.NewDevelopment()
	logger := z.Sugar()

	cfg, err := config.New()
	if err != nil || cfg == nil {
		logger.Error(err)
		panic("Cannot load config of application")
	}

	if cfg.SentryDsn != "" {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDsn,
			Environment:      cfg.AppEnv,
			Debug:            true,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			logger.Error(err)
		}
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Fatalln(err)
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	dbConnOpts, err := pq.ParseURL(cfg.DatabaseUrl)
	if err != nil {
		logger.Fatalln(err)
	}
	sqlxConn, err := sqlx.Connect("postgres", dbConnOpts)
	if err != nil {
		logger.Fatalln(err)
	}

	redisConnOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		logger.Fatalln(err)
	}
	redisClient := redis.NewClient(redisConnOpts)

	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: z}
		}),
		fx.Provide(
			func(lc fx.Lifecycle) *gorm.DB {
				lc.Append(fx.Hook{
					OnStop: func(context.Context) error {
						return d.Close()
					},
				})
				return db
			},
			func(lc fx.Lifecycle) *sqlx.DB {
				lc.Append(fx.Hook{
					OnStop: func(context.Context) error {
						return sqlxConn.Close()
					},
				})
				return sqlxConn
			},
			func(lc fx.Lifecycle) *redis.Client {
				lc.Append(fx.Hook{
					OnStop: func(context.Context) error {
						return redisClient.Close()
					},
				})
				return redisClient
			},
			func() *zap.SugaredLogger {
				return logger
			},
			func() *zap.Logger {
				return z
			},
			func() *config.Config {
				return cfg
			},
		),
		fx.Provide(middlewares.NewMiddlewares),
		fx.Provide(fiber.NewCache),
		fx.Provide(fiber.NewFiber),
		fx.Provide(v1_handlers.NewHandlers),
		fx.Invoke(routes.NewV1),
	)

	logger.Info("App started")
	app.Run()
}
