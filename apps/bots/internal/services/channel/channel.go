package channel

import (
	"log/slog"

	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Logger        *slog.Logger
	Gorm          *gorm.DB
	TwitchActions *twitchactions.TwitchActions
	WorkersPool   *workers.Pool
}

func New(opts Opts) *Service {
	return &Service{
		gorm:          opts.Gorm,
		logger:        opts.Logger,
		twitchActions: opts.TwitchActions,
		workersPool:   opts.WorkersPool,
	}
}

type Service struct {
	logger        *slog.Logger
	gorm          *gorm.DB
	twitchActions *twitchactions.TwitchActions
	workersPool   *workers.Pool
}
