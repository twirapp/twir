package app

import (
	"github.com/satont/twir/apps/ytsr/internal/grpc_impl"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"ytsr",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "ytsr"}),
	fx.Invoke(
		uptrace.NewFx("ytsr"),
		grpc_impl.New,
		func(l logger.Logger) {
			l.Info("Started")
		},
	),
)
