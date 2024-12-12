package main

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql"
	apq_cache "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/apq-cache"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/auth"
	pubclicroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/public"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/webhooks"
	"github.com/twirapp/twir/apps/api-gql/internal/minio"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	dashboard_widget_events "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	"github.com/twirapp/twir/libs/baseapp"
	commandscache "github.com/twirapp/twir/libs/cache/commands"
	keywordscacher "github.com/twirapp/twir/libs/cache/keywords"
	twitchcache "github.com/twirapp/twir/libs/cache/twitch"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"

	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "api-gql",
			},
		),
		fx.Provide(
			dashboard_widget_events.New,
			variables.New,
			timers.New,
		),
		fx.Provide(
			fx.Annotate(
				timersrepositorypgx.NewFx,
				fx.As(new(timersrepository.Repository)),
			),
		),
		fx.Provide(
			auth.NewSessions,
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
			func(config cfg.Config) events.EventsClient {
				return clients.NewEvents(config.AppEnv)
			},
			minio.New,
			twitchcache.New,
			commandscache.New,
			keywordscacher.New,
			fx.Annotate(
				wsrouter.NewNatsSubscription,
				fx.As(new(wsrouter.WsRouter)),
			),
			twir_stats.New,
			resolvers.New,
			directives.New,
			middlewares.New,
			server.New,
			apq_cache.New,
		),
		fx.Invoke(
			gql.New,
			uptrace.NewFx("api-gql"),
			pubclicroutes.New,
			webhooks.New,
			authroutes.New,
		),
	).Run()
}
