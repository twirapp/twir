package app

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/apps/dota/internal/gsi"
	"github.com/twirapp/twir/apps/dota/internal/match"
	"github.com/twirapp/twir/apps/dota/internal/processor"
	"github.com/twirapp/twir/apps/dota/internal/stats"
	"github.com/twirapp/twir/libs/baseapp"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/opendota"
	"github.com/twirapp/twir/libs/integrations/stratz"
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
		fx.Annotate(
			match.NewBusEmitter,
			fx.As(new(match.EventEmitter)),
		),
		func(config cfg.Config) *stratz.Client {
			return stratz.New(config.DotaStratzToken)
		},
		opendota.New,
		fx.Annotate(
			stats.New,
			fx.As(new(processor.WinProbabilityProvider)),
		),
		match.New,
		fx.Annotate(
			processor.New,
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
