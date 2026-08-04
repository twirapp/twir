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
)

const Service = "emotes-cacher"

var ProviderSet = wire.NewSet(
	wire.Value(baseapp.Opts{AppName: Service}),
	baseapp.NewBase,
	wire.FieldsOf(new(baseapp.Base), "Lifecycle", "Config", "Gorm", "Logger", "Bus"),
	wire.Struct(new(emotes_store.Opts), "*"),
	emotes_store.New,
	wire.Struct(new(buslistener.Opts), "*"),
	wire.Struct(new(seventv.Opts), "*"),
	wire.Struct(new(bttv.Opts), "*"),
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func NewApplication(
	lifecycle *lifecycle.Lifecycle,
	logger *slog.Logger,
	busListenerOpts buslistener.Opts,
	sevenTVOpts seventv.Opts,
	bttvOpts bttv.Opts,
) (*Application, error) {
	if err := buslistener.New(busListenerOpts); err != nil {
		return nil, fmt.Errorf("create bus listener: %w", err)
	}
	if err := seventv.New(sevenTVOpts); err != nil {
		return nil, fmt.Errorf("create 7TV service: %w", err)
	}
	if err := bttv.New(bttvOpts); err != nil {
		return nil, fmt.Errorf("create BTTV service: %w", err)
	}

	logger.Info("Emotes Cacher started")
	return &Application{lifecycle: lifecycle}, nil
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}
