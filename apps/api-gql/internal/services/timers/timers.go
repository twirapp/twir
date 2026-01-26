package timers

import (
	"errors"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/twirapp/twir/libs/audit"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/repositories/plans"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm             *gorm.DB
	AuditRecorder    audit.Recorder
	Logger           *slog.Logger
	TwirBus          *buscore.Bus
	TimersRepository timersrepository.Repository
	TrmManager       trm.Manager
	PlansRepository  plans.Repository
}

func New(opts Opts) *Service {
	return &Service{
		gorm:             opts.Gorm,
		auditRecorder:    opts.AuditRecorder,
		logger:           opts.Logger,
		twirbus:          opts.TwirBus,
		timersRepository: opts.TimersRepository,
		trmManager:       opts.TrmManager,
		plansRepository:  opts.PlansRepository,
	}
}

type Service struct {
	gorm             *gorm.DB
	logger           *slog.Logger
	auditRecorder    audit.Recorder
	twirbus          *buscore.Bus
	timersRepository timersrepository.Repository
	trmManager       trm.Manager
	plansRepository  plans.Repository
}

var ErrTimerNotFound = errors.New("timer not found")
