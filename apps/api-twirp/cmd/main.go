package main

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_protected"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_unprotected"
	"github.com/satont/tsuwari/apps/api-twirp/internal/interceptors"
	"github.com/satont/tsuwari/apps/api-twirp/internal/sessions"
	"github.com/satont/tsuwari/apps/api-twirp/internal/twirp_handlers"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"net/http"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)

	fx.New(
		fx.WithLogger(func() fxevent.Logger {
			l, _ := zap.NewDevelopment()
			return &fxevent.ZapLogger{Logger: l}
		}),
		fx.Provide(
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
			func(config *cfg.Config, lc fx.Lifecycle) *redis.Client {
				redisOpts, err := redis.ParseURL(config.RedisUrl)
				if err != nil {
					logger.Sugar().Panic(err)
				}
				client := redis.NewClient(redisOpts)
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return client.Close()
					},
				})

				return client
			},
			func(r *redis.Client) *scs.SessionManager {
				return sessions.New(r)
			},
			func(config *cfg.Config, lc fx.Lifecycle) *gorm.DB {
				db, err := gorm.Open(postgres.Open(config.DatabaseUrl), &gorm.Config{
					Logger: gormLogger.Default.LogMode(gormLogger.Silent),
				})
				if err != nil {
					logger.Sugar().Panic("failed to connect database", err)
				}
				d, _ := db.DB()
				d.SetMaxOpenConns(20)
				d.SetConnMaxIdleTime(1 * time.Minute)

				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return d.Close()
					},
				})

				return db
			},
			interceptors.New,
			impl_protected.New,
			impl_unprotected.New,
			twirp_handlers.AsHandler(twirp_handlers.NewProtected),
			twirp_handlers.AsHandler(twirp_handlers.NewUnProtected),
			fx.Annotate(
				func(handlers []twirp_handlers.IHandler) *http.ServeMux {
					mux := http.NewServeMux()
					for _, route := range handlers {
						mux.Handle(route.Pattern(), route.Handler())
					}
					return mux
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
		fx.Invoke(func(mux *http.ServeMux, sessionManager *scs.SessionManager) {
			logger.Sugar().Panic(http.ListenAndServe(":3002", sessionManager.LoadAndSave(mux)))
		}),
	).Run()
}
