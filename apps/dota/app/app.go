package app

import (
	"log/slog"

	"github.com/goforj/wire"

	"github.com/twirapp/twir/apps/dota/internal/buslistener"
	"github.com/twirapp/twir/apps/dota/internal/chatalerts"
	"github.com/twirapp/twir/apps/dota/internal/gsi"
	"github.com/twirapp/twir/apps/dota/internal/match"
	"github.com/twirapp/twir/apps/dota/internal/predictions"
	"github.com/twirapp/twir/apps/dota/internal/processor"
	"github.com/twirapp/twir/apps/dota/internal/stats"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/opendota"
	"github.com/twirapp/twir/libs/integrations/stratz"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotarepositorypgx "github.com/twirapp/twir/libs/repositories/dota/pgx"
)

const Service = "dota"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	dotarepositorypgx.NewFx,
	wire.Bind(new(dotarepository.Repository), new(*dotarepositorypgx.Pgx)),
	channelsrepositorypgx.NewFx,
	wire.Bind(new(channelsrepository.Repository), new(*channelsrepositorypgx.Pgx)),
	match.NewBusEmitter,
	wire.Bind(new(match.EventEmitter), new(*match.BusEmitter)),
	NewStratzClient,
	NewOpenDotaClient,
	stats.New,
	wire.Bind(new(processor.WinProbabilityProvider), new(*stats.Stats)),
	wire.Bind(new(buslistener.StatsProvider), new(*stats.Stats)),
	match.New,
	processor.New,
	wire.Bind(new(gsi.MatchProcessor), new(*processor.Processor)),
	gsi.New,
	buslistener.New,
	chatalerts.NewRedisCooldownStore,
	wire.Bind(new(chatalerts.CooldownStore), new(*chatalerts.RedisCooldownStore)),
	chatalerts.New,
	predictions.NewRedisPredictionStore,
	wire.Bind(new(predictions.Store), new(*predictions.RedisPredictionStore)),
	predictions.New,
	predictions.NewLifecycleWorker,
	NewApplication,
)

func NewStratzClient(config cfg.Config) *stratz.Client {
	return stratz.New(config.DotaStratzToken)
}

// wire cannot inject variadic options; the default client has none.
func NewOpenDotaClient() *opendota.Client {
	return opendota.New()
}

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lc *lifecycle.Lifecycle,
	logger *slog.Logger,
	_ *gsi.Server,
	_ *buslistener.BusListener,
	_ *chatalerts.ChatAlerts,
	_ *predictions.LifecycleWorker,
) *Application {
	logger.Info("🤖 Dota service started")

	return &Application{lifecycle: lc}
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
