package app

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/apps/dota/internal/gsi"
	"github.com/twirapp/twir/apps/dota/internal/processorstub"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/otel"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotarepositorypgx "github.com/twirapp/twir/libs/repositories/dota/pgx"
	"go.uber.org/fx"
)

var App = fx.Module(
	"dota",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "dota"}),
	fx.Provide(
		fx.Annotate(
			dotarepositorypgx.NewFx,
			fx.As(new(dotarepository.Repository)),
		),
		// TODO(task-6): replace processorstub with the real match processor.
		fx.Annotate(
			processorstub.New,
			fx.As(new(gsi.MatchProcessor)),
		),
		gsi.New,
	),
	fx.Invoke(
		otel.NewFx("dota"),
		func(s *gsi.Server, lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error { return s.Start() },
				OnStop:  func(ctx context.Context) error { return s.Stop(ctx) },
			})
		},
		func(l *slog.Logger) {
			l.Info("🤖 Dota service started")
		},
	),
)
