package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/api/internal/handlers"
	"github.com/twirapp/twir/apps/api/internal/impl_protected"
	"github.com/twirapp/twir/apps/api/internal/interceptors"
	"github.com/twirapp/twir/apps/api/internal/proxy"
	"github.com/twirapp/twir/apps/api/internal/sessions"
	"github.com/twirapp/twir/apps/api/internal/twirp_handlers"
	"github.com/twirapp/twir/libs/baseapp"
	channelseventswithoperations "github.com/twirapp/twir/libs/cache/channels_events_with_operations"
	commandswithgroupsandresponsescache "github.com/twirapp/twir/libs/cache/commands"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/websockets"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	commandswithgroupsandresponsespostgres "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/pgx"
	channelseventsrepository "github.com/twirapp/twir/libs/repositories/events"
	channelseventsrepositorypostgres "github.com/twirapp/twir/libs/repositories/events/pgx"
	"github.com/twirapp/twir/libs/otel"
	"go.uber.org/fx"
)

var App = fx.Options(
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "api"}),
	fx.Provide(
		fx.Annotate(
			channelsintegrationsspotifypgx.NewFx,
			fx.As(new(channelsintegrationsspotify.Repository)),
		),
		fx.Annotate(
			channelseventsrepositorypostgres.NewFx,
			fx.As(new(channelseventsrepository.Repository)),
		),
		fx.Annotate(
			commandswithgroupsandresponsespostgres.NewFx,
			fx.As(new(commandswithgroupsandresponsesrepository.Repository)),
		),
	),
	fx.Provide(
		func(c cfg.Config) websockets.WebsocketClient {
			return clients.NewWebsocket(c.AppEnv)
		},
		func(r *redis.Client) *scs.SessionManager {
			return sessions.New(r)
		},
		channelseventswithoperations.New,
		commandswithgroupsandresponsescache.New,
		interceptors.New,
		impl_protected.New,
		handlers.AsHandler(twirp_handlers.NewProtected),
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
		otel.NewFx("api"),
		func(
			mux *http.ServeMux,
			sessionManager *scs.SessionManager,
			l *slog.Logger,
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
