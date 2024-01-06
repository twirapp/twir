package app

import (
	"github.com/satont/twir/apps/ytsr/internal/grpc_impl"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	"go.uber.org/fx"
)

var App = fx.Module(
	"ytsr",
	fx.NopLogger,
	fx.Provide(
		cfg.NewFx,
		twirsentry.NewFx(twirsentry.NewFxOpts{Service: "ytsr"}),
		logger.NewFx(logger.Opts{Service: "ytsr"}),
	),
	fx.Invoke(
		grpc_impl.New,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)
