package main

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/bots"
	"github.com/satont/twir/apps/bots/internal/gorm"
	"github.com/satont/twir/apps/bots/internal/grpc_impl"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/pubsub"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			cfg.NewFx,
			func(config cfg.Config) (*sentry.Client, error) {
				if config.SentryDsn == "" {
					return nil, nil
				}

				s, err := sentry.NewClient(
					sentry.ClientOptions{
						Dsn:              config.SentryDsn,
						AttachStacktrace: true,
					},
				)

				return s, err
			},
			func(config cfg.Config, s *sentry.Client) logger.Logger {
				return logger.New(
					logger.Opts{
						Env:     config.AppEnv,
						Service: "timers",
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
		fx.NopLogger,
		fx.Invoke(
			func(config cfg.Config) {
				if config.AppEnv != "development" {
					http.Handle("/metrics", promhttp.Handler())
					go http.ListenAndServe("0.0.0.0:3000", nil)
				}
			},
			func(config cfg.Config) {
				if config.SentryDsn != "" {
					sentry.Init(
						sentry.ClientOptions{
							Dsn:              config.SentryDsn,
							Environment:      config.AppEnv,
							Debug:            true,
							TracesSampleRate: 1.0,
						},
					)
				}
			},
			grpc_impl.NewServer,
		),
	).Run()
}
