package main

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql"
	apq_cache "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/apq-cache"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/directives"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/resolvers"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-public/auth"
	pubclicroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http-public/public"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http-public/webhooks"
	"github.com/twirapp/twir/apps/api-gql/internal/minio"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/server/middlewares"
	admin_actions "github.com/twirapp/twir/apps/api-gql/internal/services/admin-actions"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
	dashboard_widget_events "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	twir_users "github.com/twirapp/twir/apps/api-gql/internal/services/twir-users"
	twitch_channels "github.com/twirapp/twir/apps/api-gql/internal/services/twitch-channels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
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

	variablesrepository "github.com/twirapp/twir/libs/repositories/variables"
	variablespgx "github.com/twirapp/twir/libs/repositories/variables/pgx"

	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"

	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
	keywordsrepositorypgx "github.com/twirapp/twir/libs/repositories/keywords/pgx"

	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"

	badgesrepository "github.com/twirapp/twir/libs/repositories/badges"
	badgesrepositorypgx "github.com/twirapp/twir/libs/repositories/badges/pgx"

	badgesusersrepository "github.com/twirapp/twir/libs/repositories/badges-users"
	badgesusersrepositorypgx "github.com/twirapp/twir/libs/repositories/badges-users/pgx"

	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"

	userswithchannelrepository "github.com/twirapp/twir/libs/repositories/users-with-channel"
	userswithchannelrepositorypgx "github.com/twirapp/twir/libs/repositories/users-with-channel/pgx"
)

func main() {
	fx.New(
		baseapp.CreateBaseApp(
			baseapp.Opts{
				AppName: "api-gql",
			},
		),
		fx.Provide(
			twitchcache.New,
		),
		// services
		fx.Provide(
			dashboard_widget_events.New,
			variables.New,
			timers.New,
			keywords.New,
			audit_logs.New,
			admin_actions.New,
			badges.New,
			badges_users.New,
			badges_with_users.New,
			users.New,
			twitch_channels.New,
			twir_users.New,
		),
		// repositories
		fx.Provide(
			fx.Annotate(
				timersrepositorypgx.NewFx,
				fx.As(new(timersrepository.Repository)),
			),
			fx.Annotate(
				variablespgx.NewFx,
				fx.As(new(variablesrepository.Repository)),
			),
			fx.Annotate(
				keywordsrepositorypgx.NewFx,
				fx.As(new(keywordsrepository.Repository)),
			),
			fx.Annotate(
				channelsrepositorypgx.NewFx,
				fx.As(new(channelsrepository.Repository)),
			),
			fx.Annotate(
				badgesrepositorypgx.NewFx,
				fx.As(new(badgesrepository.Repository)),
			),
			fx.Annotate(
				badgesusersrepositorypgx.NewFx,
				fx.As(new(badgesusersrepository.Repository)),
			),
			fx.Annotate(
				usersrepositorypgx.NewFx,
				fx.As(new(usersrepository.Repository)),
			),
			fx.Annotate(
				userswithchannelrepositorypgx.NewFx,
				fx.As(new(userswithchannelrepository.Repository)),
			),
		),
		// grpc clients
		fx.Provide(
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
			func(config cfg.Config) events.EventsClient {
				return clients.NewEvents(config.AppEnv)
			},
		),
		// app itself
		fx.Provide(
			auth.NewSessions,
			minio.New,
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
