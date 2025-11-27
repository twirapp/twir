package app

import (
	"log/slog"

	"github.com/twirapp/twir/apps/ytsr/internal/bus_listener"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

var App = fx.Module(
	"ytsr",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "ytsr"}),
	fx.Invoke(
		uptrace.NewFx("ytsr"),
		bus_listener.New,
		func(l *slog.Logger) {
			l.Info("Started")
		},
	),
)
