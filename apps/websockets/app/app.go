package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	bus_listener "github.com/twirapp/twir/apps/websockets/internal/bus-listener"
	"github.com/twirapp/twir/apps/websockets/internal/grpc_impl"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/alerts"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/be_right_back"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/dudes"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/obs"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/registry/overlays"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/tts"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"

	kappagenrepository "github.com/twirapp/twir/libs/repositories/overlays_kappagen"
	kappagenrepositorypgx "github.com/twirapp/twir/libs/repositories/overlays_kappagen/pgx"
)

const service = "Websockets"

var App = fx.Module(
	service,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: service}),
	fx.Provide(
		fx.Annotate(
			kappagenrepositorypgx.NewFx,
			fx.As(new(kappagenrepository.Repository)),
		),
	),
	fx.Provide(
		fx.Annotate(
			alertsrepositorypgx.NewFx,
			fx.As(new(alertsrepository.Repository)),
		),
		tts.NewTts,
		obs.NewObs,
		youtube.NewYouTube,
		alerts.NewAlerts,
		channelalertscache.New,
		overlays.New,
		be_right_back.New,
		dudes.New,
	),
	fx.Invoke(
		uptrace.NewFx(service),
		bus_listener.New,
		func() {
			http.Handle("/metrics", promhttp.Handler())

			go http.ListenAndServe(":3004", nil)
		},
		grpc_impl.NewGrpcImplementation,
		func(l logger.Logger) {
			l.Info(service + " started")
		},
	),
)
