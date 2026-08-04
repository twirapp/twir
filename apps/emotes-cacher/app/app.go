package app

import (
	"fmt"
	"log/slog"

	"github.com/goforj/wire"
	buslistener "github.com/twirapp/twir/apps/emotes-cacher/internal/bus-listener"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emotes_store"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/bttv"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"gorm.io/gorm"
)

const Service = "emotes-cacher"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.ProviderSet,
	emotes_store.New,
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	bus *buscore.Bus,
	gorm *gorm.DB,
	cfg config.Config,
	emotesStore *emotes_store.EmotesStore,
) (*Application, error) {
	if err := buslistener.New(lifecycle, logger, bus, emotesStore); err != nil {
		return nil, fmt.Errorf("create bus listener: %w", err)
	}
	if err := seventv.New(lifecycle, gorm, cfg, logger, emotesStore); err != nil {
		return nil, fmt.Errorf("create 7TV service: %w", err)
	}
	if err := bttv.New(lifecycle, gorm, logger, emotesStore); err != nil {
		return nil, fmt.Errorf("create BTTV service: %w", err)
	}

	logger.Info("Emotes Cacher started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
