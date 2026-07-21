package app

import (
	"log/slog"

	bus_listener "github.com/twirapp/twir/apps/timers/internal/bus-listener"
	"github.com/twirapp/twir/apps/timers/internal/manager"
	"github.com/twirapp/twir/libs/baseapp"
	channelcache "github.com/twirapp/twir/libs/cache/channel"
	"github.com/twirapp/twir/libs/otel"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	"github.com/twirapp/twir/libs/repositories/streams"
	streamsrepositorypostgres "github.com/twirapp/twir/libs/repositories/streams/datasource/postgres"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersrepositorypgx "github.com/twirapp/twir/libs/repositories/timers/pgx"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
)

var App = fx.Module(
	"timers",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "timers"}),
	fx.Provide(
		fx.Annotate(
			timersrepositorypgx.NewFx,
			fx.As(new(timersrepository.Repository)),
		),
		fx.Annotate(
			channelsrepositorypgx.NewFx,
			fx.As(new(channelsrepository.Repository)),
		),
		fx.Annotate(
			streamsrepositorypostgres.NewFx,
			fx.As(new(streams.Repository)),
		),
		channelservice.NewChannelService,
		channelcache.New,
		manager.New,
	),
	fx.Invoke(
		otel.NewFx("timers"),
		bus_listener.New,
		func(l *slog.Logger) {
			l.Info("🚀 Timers service started")
		},
	),
)
