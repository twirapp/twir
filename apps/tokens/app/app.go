package app

import (
	"github.com/satont/twir/apps/tokens/internal/grpc_impl"
	"github.com/satont/twir/apps/tokens/internal/redis"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"tokens",
	baseapp.CreateBaseApp("tokens"),
	fx.Provide(
		redis.NewRedisLock,
	),
	fx.Invoke(
		uptrace.NewFx("tokens"),
		grpc_impl.NewTokens,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)
