package app

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	bus_listener "github.com/twirapp/twir/apps/websockets/internal/bus-listener"
	"github.com/twirapp/twir/apps/websockets/internal/grpc_impl"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/alerts"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/dudes"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/registry/overlays"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/twirapp/twir/libs/baseapp"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	"github.com/twirapp/twir/libs/repositories/channels_overlays"
	channelsoverlayspgx "github.com/twirapp/twir/libs/repositories/channels_overlays/pgx"
	"github.com/twirapp/twir/libs/otel"
	"github.com/twirapp/twir/libs/wsrouter"
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
		fx.Annotate(
			wsrouter.NewNatsWsRouterFx,
			fx.As(new(wsrouter.WsRouter)),
		),
		fx.Annotate(
			channelsoverlayspgx.NewFx,
			fx.As(new(channels_overlays.Repository)),
		),
		youtube.NewYouTube,
		alerts.NewAlerts,
		channelalertscache.New,
		overlays.New,
		dudes.New,
	),
	fx.Invoke(
		otel.NewFx(service),
		bus_listener.New,
		func() {
			http.Handle("/metrics", promhttp.Handler())

			go http.ListenAndServe(":3004", nil)
		},
		grpc_impl.NewGrpcImplementation,
		func(l *slog.Logger) {
			l.Info(service + " started")
		},
	),
)
