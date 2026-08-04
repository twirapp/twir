package timers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/twirapp/twir/libs/audit"
	buscore "github.com/twirapp/twir/libs/bus-core"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
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
		timerLifecycle:   busTimerLifecycle{bus: opts.TwirBus},
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
	timerLifecycle   timerLifecycle
	timersRepository timersrepository.Repository
	trmManager       trm.Manager
	plansRepository  plans.Repository
}

type timerLifecycle interface {
	Publish(context.Context, bool, timersbusservice.AddOrRemoveTimerRequest) error
}

type busTimerLifecycle struct {
	bus *buscore.Bus
}

func (l busTimerLifecycle) Publish(
	ctx context.Context,
	enabled bool,
	request timersbusservice.AddOrRemoveTimerRequest,
) error {
	if enabled {
		return l.bus.Timers.AddTimer.Publish(ctx, request)
	}

	return l.bus.Timers.RemoveTimer.Publish(ctx, request)
}

var ErrTimerNotFound = errors.New("timer not found")
