package main

import (
	"context"
	"net/http"
	"time"

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
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	fx.New(
		fx.WithLogger(
			func() fxevent.Logger {
				l, _ := zap.NewDevelopment()
				return &fxevent.ZapLogger{Logger: l}
			},
		),
		fx.Provide(
			func() *zap.Logger {
				return logger
			},
			func() *cfg.Config {
				config, err := cfg.New()
				if err != nil {
					panic(err)
				}
				return config
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
			func(config *cfg.Config, lc fx.Lifecycle) *redis.Client {
				redisOpts, err := redis.ParseURL(config.RedisUrl)
				if err != nil {
					logger.Sugar().Panic(err)
				}
				client := redis.NewClient(redisOpts)
				lc.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							return client.Close()
						},
					},
				)

				return client
			},
			func(r *redis.Client) *scs.SessionManager {
				return sessions.New(r)
			},
			func(config *cfg.Config, lc fx.Lifecycle) *gorm.DB {
				db, err := gorm.Open(
					postgres.Open(config.DatabaseUrl), &gorm.Config{
						Logger: gormLogger.Default.LogMode(gormLogger.Silent),
					},
				)
				if err != nil {
					logger.Sugar().Panic("failed to connect database", err)
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

				return db
			},
			interceptors.New,
			impl_protected.New,
			impl_unprotected.New,
			handlers.AsHandler(twirp_handlers.NewProtected),
			handlers.AsHandler(twirp_handlers.NewUnProtected),
			handlers.AsHandler(webhooks.NewDonateStream),
			handlers.AsHandler(webhooks.NewDonatello),
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
			func(mux *http.ServeMux, sessionManager *scs.SessionManager) {
				logger.Sugar().Info("Api started")
				logger.Sugar().Panic(http.ListenAndServe(":3002", sessionManager.LoadAndSave(mux)))
			},
		),
	).Run()
}
