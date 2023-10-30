package main

import (
	"net/http"

	"github.com/satont/twir/apps/websockets/internal/gorm"
	"github.com/satont/twir/apps/websockets/internal/grpc_impl"
	"github.com/satont/twir/apps/websockets/internal/namespaces/alerts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/chat"
	"github.com/satont/twir/apps/websockets/internal/namespaces/obs"
	"github.com/satont/twir/apps/websockets/internal/namespaces/registry/overlays"
	"github.com/satont/twir/apps/websockets/internal/namespaces/tts"
	"github.com/satont/twir/apps/websockets/internal/namespaces/youtube"
	"github.com/satont/twir/apps/websockets/internal/redis"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	"go.uber.org/fx"
)

func main() {
	const service = "Websockets"

	fx.New(
		fx.NopLogger,
		fx.Provide(
			config.NewFx,
			twirsentry.NewFx(twirsentry.NewFxOpts{Service: service}),
			logger.NewFx(logger.Opts{Service: service}),
			redis.New,
			gorm.New,
			func(cfg config.Config) bots.BotsClient {
				return clients.NewBots(cfg.AppEnv)
			},
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
			chat.New,
			overlays.New,
		),
		fx.Invoke(
			func() {
				go http.ListenAndServe(":3004", nil)
			},
			grpc_impl.NewGrpcImplementation,
			func(l logger.Logger) {
				l.Info(service + " started")
			},
		),
	).Run()

}
