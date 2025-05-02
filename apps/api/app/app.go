package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/api/internal/files"
	"github.com/satont/twir/apps/api/internal/handlers"
	"github.com/satont/twir/apps/api/internal/impl_protected"
	"github.com/satont/twir/apps/api/internal/impl_unprotected"
	"github.com/satont/twir/apps/api/internal/interceptors"
	"github.com/satont/twir/apps/api/internal/proxy"
	"github.com/satont/twir/apps/api/internal/sessions"
	"github.com/satont/twir/apps/api/internal/twirp_handlers"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	ttscache "github.com/twirapp/twir/libs/cache/tts"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/discord"
	"github.com/twirapp/twir/libs/grpc/integrations"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Options(
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "api"}),
	fx.Provide(
		fx.Annotate(
			channelsintegrationsspotifypgx.NewFx,
			fx.As(new(channelsintegrationsspotify.Repository)),
		),
	),
	fx.Provide(
		func(c cfg.Config) tokens.TokensClient {
			return clients.NewTokens(c.AppEnv)
		},
		func(c cfg.Config) integrations.IntegrationsClient {
			return clients.NewIntegrations(c.AppEnv)
		},
		func(c cfg.Config) parser.ParserClient {
			return clients.NewParser(c.AppEnv)
		},
		func(c cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(c.AppEnv)
		},
		func(c cfg.Config) discord.DiscordClient {
			return clients.NewDiscord(c.AppEnv)
		},
		func(r *redis.Client) *scs.SessionManager {
			return sessions.New(r)
		},
		interceptors.New,
		impl_protected.New,
		impl_unprotected.New,
		ttscache.NewTTSSettings,
		handlers.AsHandler(twirp_handlers.NewProtected),
		handlers.AsHandler(twirp_handlers.NewUnProtected),
		handlers.AsHandler(files.NewFiles),
		handlers.AsHandler(proxy.New),
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
	fx.Invoke(
		uptrace.NewFx("api"),
		func(
			mux *http.ServeMux,
			sessionManager *scs.SessionManager,
			l logger.Logger,
			lc fx.Lifecycle,
		) error {
			server := &http.Server{
				Addr:    "0.0.0.0:3002",
				Handler: sessionManager.LoadAndSave(mux),
			}

			lc.Append(
				fx.Hook{
					OnStart: func(_ context.Context) error {
						go func() {
							l.Info("Started", slog.String("port", "3002"))
							err := server.ListenAndServe()
							if err != nil && !errors.Is(err, http.ErrServerClosed) {
								panic(err)
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {
						return server.Shutdown(ctx)
					},
				},
			)

			return nil
		},
	),
)
