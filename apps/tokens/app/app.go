package app

import (
	"log/slog"

	"github.com/twirapp/twir/apps/tokens/internal/bus_listener"
	"github.com/twirapp/twir/apps/tokens/internal/redis"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/otel"
	"go.uber.org/fx"

	userswithtokensrepository "github.com/twirapp/twir/libs/repositories/userswithtoken"
	userswithtokensrepositorypgx "github.com/twirapp/twir/libs/repositories/userswithtoken/pgx"

	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokensrepositorypostgres "github.com/twirapp/twir/libs/repositories/tokens/datasources/postgres"

	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	kickbotsrepositorypgx "github.com/twirapp/twir/libs/repositories/kick_bots/pgx"
	channelsintegrationsrepository "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_integrations/datasource/postgres"
	channelsintegrationsspotifyrepository "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	channelsintegrationsspotifyrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/pgx"
	integrationsrepository "github.com/twirapp/twir/libs/repositories/integrations"
	integrationsrepositorypgx "github.com/twirapp/twir/libs/repositories/integrations/datasource/postgres"
)

var App = fx.Module(
	"tokens",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "tokens"}),
	fx.Provide(
		fx.Annotate(
			userswithtokensrepositorypgx.NewFx,
			fx.As(new(userswithtokensrepository.Repository)),
		),
		fx.Annotate(
			tokensrepositorypostgres.NewFx,
			fx.As(new(tokensrepository.Repository)),
		),
		fx.Annotate(
			usersrepositorypgx.NewFx,
			fx.As(new(usersrepository.Repository)),
		),
		fx.Annotate(
			kickbotsrepositorypgx.NewFx,
			fx.As(new(kickbotsrepository.Repository)),
		),
		fx.Annotate(
			channelsintegrationsrepositorypgx.NewFx,
			fx.As(new(channelsintegrationsrepository.Repository)),
		),
		fx.Annotate(
			channelsintegrationsspotifyrepositorypgx.NewFx,
			fx.As(new(channelsintegrationsspotifyrepository.Repository)),
		),
		fx.Annotate(
			integrationsrepositorypgx.NewFx,
			fx.As(new(integrationsrepository.Repository)),
		),
		redis.NewRedisLock,
	),
	fx.Invoke(
		otel.NewFx("tokens"),
		bus_listener.NewTokens,
		func(l *slog.Logger) {
			l.Info("Started")
		},
	),
)
