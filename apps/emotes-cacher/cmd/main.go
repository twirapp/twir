package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"github.com/satont/twir/apps/emotes-cacher/internal/di"
	"github.com/satont/twir/apps/emotes-cacher/internal/grpc_impl"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/emotes_cacher"
	"github.com/satont/twir/libs/grpc/servers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	do.ProvideValue[config.Config](di.Provider, *cfg)

	if cfg.SentryDsn != "" {
		sentry.Init(
			sentry.ClientOptions{
				Dsn:              cfg.SentryDsn,
				Environment:      cfg.AppEnv,
				Debug:            true,
				TracesSampleRate: 1.0,
			},
		)
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	do.ProvideValue[zap.Logger](di.Provider, *logger)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	do.ProvideValue[gorm.DB](di.Provider, *db)

	redisUrl, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisUrl)
	do.ProvideValue[redis.Client](di.Provider, *redisClient)

	logger.Info("Emotes cacher microservice started")

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.EMOTES_CACHER_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	emotes_cacher.RegisterEmotesCacherServer(grpcServer, grpc_impl.NewEmotesCacher())
	go grpcServer.Serve(lis)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
	grpcServer.Stop()
}
