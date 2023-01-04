package main

import (
	"fmt"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/dota"
	"github.com/satont/tsuwari/libs/grpc/generated/eval"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"

	"github.com/satont/tsuwari/apps/parser/internal/commands"
	myRedis "github.com/satont/tsuwari/apps/parser/internal/config/redis"
	"github.com/satont/tsuwari/apps/parser/internal/variables"

	config "github.com/satont/tsuwari/libs/config"

	twitch "github.com/satont/tsuwari/apps/parser/internal/config/twitch"
	"github.com/satont/tsuwari/apps/parser/internal/twitch/user"

	"github.com/getsentry/sentry-go"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/satont/tsuwari/apps/parser/internal/grpc_impl"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/servers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v9"
)

func main() {
	cfg, err := config.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	do.ProvideValue[config.Config](di.Provider, *cfg)

	if cfg.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	if cfg.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDsn,
			Environment:      cfg.AppEnv,
			Debug:            true,
			TracesSampleRate: 1.0,
		})
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	do.ProvideValue[gorm.DB](di.Provider, *db)

	r := myRedis.New(cfg.RedisUrl)
	defer r.Close()

	do.ProvideValue[redis.Client](di.Provider, *r)

	do.ProvideValue[websockets.WebsocketClient](di.Provider, clients.NewWebsocket(cfg.AppEnv))
	do.ProvideValue[bots.BotsClient](di.Provider, clients.NewBots(cfg.AppEnv))
	do.ProvideValue[dota.DotaClient](di.Provider, clients.NewDota(cfg.AppEnv))
	do.ProvideValue[eval.EvalClient](di.Provider, clients.NewEval(cfg.AppEnv))

	do.ProvideValue[users_twitch_auth.UsersTokensService](di.Provider, *users_twitch_auth.New())
	do.ProvideValue[twitch.Twitch](di.Provider, *twitch.New(*cfg))

	do.ProvideValue[variables.Variables](di.Provider, variables.New())

	do.ProvideValue[commands.Commands](di.Provider, commands.New())

	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.PARSER_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	parser.RegisterParserServer(grpcServer, grpc_impl.NewServer())
	go grpcServer.Serve(lis)

	logger.Info("Started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	log.Fatalf("Exiting")
	grpcServer.Stop()
	d.Close()
}
