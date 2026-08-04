package app

import (
	"log/slog"
	"net/http"

	"github.com/goforj/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	buslistener "github.com/twirapp/twir/apps/websockets/internal/bus-listener"
	"github.com/twirapp/twir/apps/websockets/internal/grpc_impl"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/alerts"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/dudes"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/registry/overlays"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	channelalertscache "github.com/twirapp/twir/libs/cache/channel_alerts"
	grpcwebsockets "github.com/twirapp/twir/libs/grpc/websockets"
	alertsrepository "github.com/twirapp/twir/libs/repositories/alerts"
	alertsrepositorypgx "github.com/twirapp/twir/libs/repositories/alerts/pgx"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	channelsoverlaysrepository "github.com/twirapp/twir/libs/repositories/channels_overlays"
	channelsoverlaysrepositorypgx "github.com/twirapp/twir/libs/repositories/channels_overlays/pgx"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersrepositorypgx "github.com/twirapp/twir/libs/repositories/users/pgx"
	"github.com/twirapp/twir/libs/wsrouter"
)

const Service = "Websockets"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	alertsrepositorypgx.NewFx,
	wire.Bind(new(alertsrepository.Repository), new(*alertsrepositorypgx.Pgx)),
	channelsoverlaysrepositorypgx.NewFx,
	wire.Bind(
		new(channelsoverlaysrepository.Repository),
		new(*channelsoverlaysrepositorypgx.Pgx),
	),
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	usersrepositorypgx.NewFx,
	wire.Bind(new(usersrepository.Repository), new(*usersrepositorypgx.Pgx)),
	wsrouter.NewNatsWsRouter,
	wire.Bind(new(wsrouter.WsRouter), new(*wsrouter.WsRouterNats)),
	channelalertscache.New,
	wire.Struct(new(alerts.Opts), "*"),
	alerts.NewAlerts,
	wire.Struct(new(overlays.Opts), "*"),
	overlays.New,
	wire.Struct(new(dudes.Opts), "*"),
	dudes.New,
	wire.Struct(new(buslistener.Opts), "*"),
	buslistener.New,
	wire.Struct(new(grpc_impl.GrpcOpts), "*"),
	grpc_impl.NewGrpcImplementation,
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	_ *buslistener.BusListener,
	_ grpcwebsockets.WebsocketServer,
) *Application {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":3004", nil); err != nil {
			logger.Error("metrics server stopped", "error", err)
		}
	}()

	logger.Info(Service + " started")
	return &Application{lifecycle: lifecycle}
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
