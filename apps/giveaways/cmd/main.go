package main

import (
	"context"
	"fmt"
	"github.com/satont/tsuwari/apps/giveaways/grpc_impl"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/giveaways"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	_, appCtxCancel := context.WithCancel(context.Background())

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

	redisConnOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		logger.Fatalln(err)
	}
	redisClient := redis.NewClient(redisConnOpts)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.GIVEAWAYS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	giveaways.RegisterGiveawaysServer(grpcServer, grpc_impl.NewServer())
	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

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
			func() *config.Config { return cfg },
		),
	)

	app.Run()
	logger.Info("Giveaways microservice started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Info("Exiting")
	appCtxCancel()

}
