package app

import (
	"log/slog"

	bus_listener "github.com/twirapp/twir/apps/emotes-cacher/internal/bus-listener"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emotes_store"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/bttv"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/otel"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	"go.uber.org/fx"
)

const service = "emotes-cacher"

var App = fx.Module(
	service,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: service}),
	fx.Provide(
		emotes_store.New,
		fx.Annotate(
			usersrepositorypgx.NewFx,
			fx.As(new(usersrepository.Repository)),
		),
	),
	fx.Invoke(
		otel.NewFx("emotes-cacher"),
		bus_listener.New,
		func(l *slog.Logger) {
			l.Info("Emotes Cacher started")
		},
		seventv.New,
		bttv.New,
	),
)
