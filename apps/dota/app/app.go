package app

import (
	"log/slog"

	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/otel"
	"go.uber.org/fx"
)

var App = fx.Module(
	"dota",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "dota"}),
	fx.Invoke(
		otel.NewFx("dota"),
		func(l *slog.Logger) {
			l.Info("🤖 Dota service started")
		},
	),
)
