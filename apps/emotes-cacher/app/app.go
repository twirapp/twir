package app

import (
	bus_listener "github.com/satont/twir/apps/emotes-cacher/internal/bus-listener"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

const service = "emotes-cacher"

var App = fx.Module(
	service,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: service}),
	fx.Provide(
		buscore.NewNatsBusFx("emotes-cacher"),
	),
	fx.Invoke(
		uptrace.NewFx("emotes-cacher"),
		bus_listener.New,
		func(l logger.Logger) {
			l.Info("Emotes Cacher started")
		},
	),
)
