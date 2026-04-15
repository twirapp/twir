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
