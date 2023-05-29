package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/lib/pq"
	config "github.com/satont/tsuwari/libs/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	z, _ := zap.NewDevelopment()
	logger := z.Sugar()

	cfg, err := config.New()
	if err != nil {
		logger.Error(err)
		panic("cannot load config of application")
	}
	if err != nil {
		logger.Error(err)
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

	redisConnOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		logger.Fatalln(err)
	}
	redisClient := redis.NewClient(redisConnOpts)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", services.GIVEAWAYS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	


	app := fx.New(
		fx.WithLogger(func() fxevent.Logger { return &fxevent.ZapLogger{Logger: z} }),
		fx.Provide(
			func(lc fx.Lifecycle) *gorm.DB {
				lc.Append(fx.Hook{
					OnStop: func(context.Context) error { return d.Close() },
				})
				return db
			},
			func(lc fx.Lifecycle) *redis.Client {
				lc.Append(fx.Hook{
					OnStop: func(context.Context) error { return redisClient.Close() },
				})
				return redisClient
			},
			func() *zap.SugaredLogger { return logger },
			func() *zap.Logger { return z },
			func() *config.Config { return cfg },
		),
		fx.Provide()
	)

	logger.Info("App started")
	app.Run()
}
