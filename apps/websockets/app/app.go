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
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/grpc/clients"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

const service = "Websockets"

var App = fx.Module(
	service,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: service}),
	fx.Provide(
		func(cfg config.Config) parser.ParserClient {
			return clients.NewParser(cfg.AppEnv)
		},
		func(cfg config.Config) tokens.TokensClient {
			return clients.NewTokens(cfg.AppEnv)
		},
		tts.NewTts,
		obs.NewObs,
		youtube.NewYouTube,
		alerts.NewAlerts,
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
