package app

import (
	"github.com/satont/twir/apps/tokens/internal/bus_listener"
	"github.com/satont/twir/apps/tokens/internal/redis"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"

	userswithtokensrepository "github.com/twirapp/twir/libs/repositories/userswithtoken"
	userswithtokensrepositorypgx "github.com/twirapp/twir/libs/repositories/userswithtoken/pgx"
)

var App = fx.Module(
	"tokens",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "tokens"}),
	fx.Provide(
		fx.Annotate(
			userswithtokensrepositorypgx.NewFx,
			fx.As(new(userswithtokensrepository.Repository)),
		),
		redis.NewRedisLock,
	),
	fx.Invoke(
		uptrace.NewFx("tokens"),
		bus_listener.NewTokens,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)
