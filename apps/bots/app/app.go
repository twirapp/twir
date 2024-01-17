package app

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/bots"
	"github.com/satont/twir/apps/bots/internal/gorm"
	"github.com/satont/twir/apps/bots/internal/grpc_impl"
	"github.com/satont/twir/apps/bots/internal/pubsub_handlers"
	"github.com/satont/twir/apps/bots/pkg/tlds"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"
	sentryInternal "github.com/satont/twir/libs/sentry"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
)

var App = fx.Module(
	"bots",
	fx.Provide(
		cfg.NewFx,
		tlds.New,
		sentryInternal.NewFx(sentryInternal.NewFxOpts{Service: "bots"}),
		func(config cfg.Config, s *sentry.Client) logger.Logger {
			return logger.New(
				logger.Opts{
					Env:     config.AppEnv,
					Service: "bots",
					Sentry:  s,
				},
			)
		},
		gorm.New,
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
		bots.NewBotsService,
	),
	fx.Invoke(
		func(config cfg.Config) {
			if config.AppEnv != "development" {
				http.Handle("/metrics", promhttp.Handler())
				go http.ListenAndServe("0.0.0.0:3000", nil)
			}
		},
		pubsub_handlers.New,
		grpc_impl.NewServer,
	),
)
