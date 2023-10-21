package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/satont/twir/apps/websockets/internal/namespaces/alerts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/registry/overlays"
	"github.com/satont/twir/apps/websockets/internal/namespaces/tts"

	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/websockets/internal/grpc_impl"
	"github.com/satont/twir/apps/websockets/internal/namespaces/obs"
	"github.com/satont/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/satont/twir/apps/websockets/types"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// _, cancelAppCtx := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

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

	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseUrl), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(2)
	d.SetConnMaxIdleTime(1 * time.Minute)

	redisUrl, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(redisUrl)

	services := &types.Services{
		Gorm:   db,
		Logger: logger.Sugar(),
		Redis:  rdb,
		Grpc: &types.GrpcClients{
			Bots:   clients.NewBots(cfg.AppEnv),
			Parser: clients.NewParser(cfg.AppEnv),
		},
	}

	ttsNamespace := tts.NewTts(services)
	http.HandleFunc("/tts", ttsNamespace.HandleRequest)

	obsNamespace := obs.NewObs(services)
	http.HandleFunc("/obs", obsNamespace.HandleRequest)

	youTubeNameSpace := youtube.NewYouTube(services)
	http.HandleFunc("/youtube", youTubeNameSpace.HandleRequest)

	alertsNameSpace := alerts.NewAlerts(services)
	http.HandleFunc("/alerts", alertsNameSpace.HandleRequest)

	overlaysRegistry := overlays.New(services)
	http.HandleFunc("/registry/overlays", overlaysRegistry.HandleRequest)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.WEBSOCKET_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	websockets.RegisterWebsocketServer(
		grpcServer, grpc_impl.NewGrpcImplementation(
			&grpc_impl.GrpcOpts{
				Services: services,
				Sockets: &grpc_impl.Sockets{
					TTS:              ttsNamespace,
					OBS:              obsNamespace,
					YouTube:          youTubeNameSpace,
					Alerts:           alertsNameSpace,
					OverlaysRegistry: overlaysRegistry,
				},
			},
		),
	)

	go http.ListenAndServe(":3004", nil)
	go grpcServer.Serve(lis)

	logger.Sugar().Info("Websockets server started")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	// cancelAppCtx()
	grpcServer.GracefulStop()
	lis.Close()
	d.Close()
}
