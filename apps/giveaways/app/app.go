package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	bus_listener "github.com/twirapp/twir/apps/giveaways/internal/bus-listener"
	"github.com/twirapp/twir/apps/giveaways/internal/services"
	"github.com/twirapp/twir/libs/baseapp"
	giveawayscacher "github.com/twirapp/twir/libs/cache/giveaways"
	giveawaysrepository "github.com/twirapp/twir/libs/repositories/giveaways"
	giveawaysrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways/pgx"
	giveawaysparticipantsrepository "github.com/twirapp/twir/libs/repositories/giveaways_participants"
	giveawaysparticipantsrepositorypgx "github.com/twirapp/twir/libs/repositories/giveaways_participants/pgx"
	"go.uber.org/fx"
)

var App = fx.Module(
	"giveaways",
	baseapp.CreateBaseApp(baseapp.Opts{AppName: "giveaways"}),
	fx.Provide(
		fx.Annotate(
			giveawaysrepositorypgx.NewFx,
			fx.As(new(giveawaysrepository.Repository)),
		),
		fx.Annotate(
			giveawaysparticipantsrepositorypgx.NewFx,
			fx.As(new(giveawaysparticipantsrepository.Repository)),
		),
	),
	fx.Provide(
		giveawayscacher.New,
		services.New,
	),
	fx.Invoke(
		func(config cfg.Config) {
			if config.AppEnv != "development" {
				http.Handle("/metrics", promhttp.Handler())
				go http.ListenAndServe("0.0.0.0:3000", nil)
			}
		},
		bus_listener.New,
		func(l logger.Logger) {
			l.Info("Giveaways started")
		},
	),
)
