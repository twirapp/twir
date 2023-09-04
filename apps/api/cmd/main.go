package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/satont/twir/apps/api/internal/files"
	"github.com/satont/twir/libs/grpc/generated/eventsub"
	"github.com/satont/twir/libs/logger"

	"github.com/satont/twir/libs/grpc/generated/scheduler"
	"github.com/satont/twir/libs/grpc/generated/timers"

	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/handlers"
	"github.com/satont/twir/apps/api/internal/impl_protected"
	"github.com/satont/twir/apps/api/internal/impl_unprotected"
	"github.com/satont/twir/apps/api/internal/interceptors"
	"github.com/satont/twir/apps/api/internal/sessions"
	"github.com/satont/twir/apps/api/internal/twirp_handlers"
	"github.com/satont/twir/apps/api/internal/webhooks"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/integrations"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			func() *cfg.Config {
				config, err := cfg.New()
				if err != nil {
					panic(err)
				}
				return config
			},
			func(config *cfg.Config) (*sentry.Client, error) {
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
			func(config *cfg.Config, s *sentry.Client) logger.Logger {
				return logger.New(
					logger.Opts{
						Env:     config.AppEnv,
						Service: "api",
						Sentry:  s,
					},
				)
			},
			func(c *cfg.Config) tokens.TokensClient {
				return clients.NewTokens(c.AppEnv)
			},
			func(c *cfg.Config) bots.BotsClient {
				return clients.NewBots(c.AppEnv)
			},
			func(c *cfg.Config) integrations.IntegrationsClient {
				return clients.NewIntegrations(c.AppEnv)
			},
			func(c *cfg.Config) parser.ParserClient {
				return clients.NewParser(c.AppEnv)
			},
			func(c *cfg.Config) events.EventsClient {
				return clients.NewEvents(c.AppEnv)
			},
			func(c *cfg.Config) websockets.WebsocketClient {
				return clients.NewWebsocket(c.AppEnv)
			},
			func(c *cfg.Config) scheduler.SchedulerClient {
				return clients.NewScheduler(c.AppEnv)
			},
			func(c *cfg.Config) timers.TimersClient {
				return clients.NewTimers(c.AppEnv)
			},
			func(c *cfg.Config) eventsub.EventSubClient {
				return clients.NewEventSub(c.AppEnv)
			},
			func(config *cfg.Config, lc fx.Lifecycle) (*redis.Client, error) {
				redisOpts, err := redis.ParseURL(config.RedisUrl)
				if err != nil {
					return nil, err
				}
				client := redis.NewClient(redisOpts)
				lc.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							return client.Close()
						},
					},
				)

				return client, nil
			},
			func(r *redis.Client) *scs.SessionManager {
				return sessions.New(r)
			},
			func(config *cfg.Config, lc fx.Lifecycle) (*gorm.DB, error) {
				db, err := gorm.Open(
					postgres.Open(config.DatabaseUrl),
				)
				if err != nil {
					return nil, err
				}
				d, _ := db.DB()
				d.SetMaxOpenConns(5)
				d.SetConnMaxIdleTime(1 * time.Minute)

				lc.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							return d.Close()
						},
					},
				)

				return db, nil
			},
			interceptors.New,
			impl_protected.New,
			impl_unprotected.New,
			handlers.AsHandler(twirp_handlers.NewProtected),
			handlers.AsHandler(twirp_handlers.NewUnProtected),
			handlers.AsHandler(webhooks.NewDonateStream),
			handlers.AsHandler(webhooks.NewDonatello),
			handlers.AsHandler(files.NewFiles),
			fx.Annotate(
				func(handlers []handlers.IHandler) *http.ServeMux {
					mux := http.NewServeMux()
					for _, route := range handlers {
						mux.Handle(route.Pattern(), route.Handler())
					}
					return mux
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
		fx.NopLogger,
		fx.Invoke(
			func(mux *http.ServeMux, sessionManager *scs.SessionManager, l logger.Logger) {
				l.Info("Started")
				l.Error(
					"Crashed",
					slog.Any("err", http.ListenAndServe("0.0.0.0:3002", sessionManager.LoadAndSave(mux))),
				)
			},
		),
	).Run()
}
