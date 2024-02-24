package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/gorm"
	"github.com/satont/twir/apps/bots/internal/grpc"
	"github.com/satont/twir/apps/bots/internal/messagehandler"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	"github.com/satont/twir/apps/bots/internal/nats"
	"github.com/satont/twir/apps/bots/internal/pubsub_handlers"
	"github.com/satont/twir/apps/bots/internal/queuelistener"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/satont/twir/apps/bots/pkg/tlds"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"
	"github.com/satont/twir/libs/sentry"
	"github.com/satont/twir/libs/types/types/services"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"bots",
	fx.Provide(
		cfg.NewFx,
		tlds.New,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: "bots"}),
		logger.NewFx(logger.Opts{Service: "bots"}),
		gorm.New,
		nats.New,
		uptrace.NewFx("bots"),
		services.NewNatsBus,
		func(config cfg.Config) (*pubsub.PubSub, error) {
			return pubsub.NewPubSub(config.RedisUrl)
		},
		func(config cfg.Config) tokens.TokensClient {
			return clients.NewTokens(config.AppEnv)
		},
		func(config cfg.Config) events.EventsClient {
			return clients.NewEvents(config.AppEnv)
		},
		func(config cfg.Config) parser.ParserClient {
			return clients.NewParser(config.AppEnv)
		},
		func(config cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(config.AppEnv)
		},
		func(config cfg.Config) (*redis.Client, error) {
			redisOpts, err := redis.ParseURL(config.RedisUrl)
			if err != nil {
				return nil, err
			}

			return redis.NewClient(redisOpts), nil
		},
		twitchactions.New,
		moderationhelpers.New,
		messagehandler.New,
		queuelistener.New,
	),
	fx.Invoke(
		uptrace.NewFx("bots"),
		nats.New,
		queuelistener.New,
		func(config cfg.Config) {
			if config.AppEnv != "development" {
				http.Handle("/metrics", promhttp.Handler())
				go http.ListenAndServe("0.0.0.0:3000", nil)
			}
		},
		pubsub_handlers.New,
		grpc.New,
		func(l logger.Logger) {
			l.Info("Bots started")
		},
	),
)
