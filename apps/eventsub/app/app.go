package app

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/eventsub/internal/grpc"
	"github.com/satont/twir/apps/eventsub/internal/handler"
	"github.com/satont/twir/apps/eventsub/internal/manager"
	"github.com/satont/twir/apps/eventsub/internal/pubsub"
	"github.com/satont/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/grpc/bots"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var App = fx.Options(
	fx.Provide(
		cfg.NewFx,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: "eventsub"}),
		logger.NewFx(logger.Opts{Service: "eventsub"}),
		uptrace.NewFx("eventsub"),
		func(config cfg.Config) (*gorm.DB, error) {
			db, err := gorm.Open(
				postgres.Open(config.DatabaseUrl), &gorm.Config{
					Logger: gormLogger.Default.LogMode(gormLogger.Silent),
				},
			)
			if err != nil {
				return nil, err
			}
			d, _ := db.DB()
			d.SetMaxIdleConns(1)
			d.SetMaxOpenConns(10)
			d.SetConnMaxLifetime(time.Hour)
			return db, nil
		},
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		func(config cfg.Config) events.EventsClient {
			return clients.NewEvents(config.AppEnv)
		},
		func(config cfg.Config) bots.BotsClient {
			return clients.NewBots(config.AppEnv)
		},
		func(config cfg.Config) parser.ParserClient {
			return clients.NewParser(config.AppEnv)
		},
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		func(config cfg.Config) (*redis.Client, error) {
			redisUrl, err := redis.ParseURL(config.RedisUrl)
			if err != nil {
				return nil, err
			}

			redisClient := redis.NewClient(redisUrl)
			return redisClient, nil
		},
		buscore.NewNatsBusFx("eventsub"),
		pubsub.New,
		tunnel.New,
		manager.NewCreds,
		manager.NewManager,
		handler.New,
	),
	fx.Invoke(
		uptrace.NewFx("eventsub"),
		handler.New,
		grpc.New,
	),
)
