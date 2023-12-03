package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satont/twir/libs/logger"

	"github.com/getsentry/sentry-go"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/samber/do"
	"github.com/satont/twir/apps/tokens/internal/di"
	"github.com/satont/twir/apps/tokens/internal/grpc_impl"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
				Debug:            false,
				TracesSampleRate: 1.0,
				AttachStacktrace: true,
			},
		)
	}

	l := logger.New(
		logger.Opts{
			Env:     cfg.AppEnv,
			Service: "tokens",
		},
	)

	do.ProvideValue[logger.Logger](di.Provider, l)

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(2)
	d.SetConnMaxIdleTime(1 * time.Minute)

	do.ProvideValue[gorm.DB](di.Provider, *db)

	// redis lock
	redisUrl, err := goredislib.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := goredislib.NewClient(redisUrl)
	pool := goredis.NewPool(redisClient)
	do.ProvideValue[redsync.Redsync](di.Provider, *redsync.New(pool))
	//

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.TOKENS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	tokens.RegisterTokensServer(grpcServer, grpc_impl.NewTokens())
	go grpcServer.Serve(lis)

	l.Info("Started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
	grpcServer.Stop()
}
