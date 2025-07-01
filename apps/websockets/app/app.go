package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	bus_listener "github.com/satont/twir/apps/websockets/internal/bus-listener"
	"github.com/satont/twir/apps/websockets/internal/grpc_impl"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/alerts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/be_right_back"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/dudes"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/kappagen"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/obs"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/registry/overlays"
	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/tts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

const service = "Websockets"

var App = fx.Module(
	service,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: service}),
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
		kappagen.New,
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
