package timers

import (
	"errors"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/twirapp/twir/libs/audit"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/repositories/plans"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	auditRecorder audit.Recorder,
	logger *slog.Logger,
	twirBus *buscore.Bus,
	timersRepository timersrepository.Repository,
	trmManager trm.Manager,
	plansRepository plans.Repository,
) *Service {
	return &Service{
		gorm:             db,
		auditRecorder:    auditRecorder,
		logger:           logger,
		twirbus:          twirBus,
		timersRepository: timersRepository,
		trmManager:       trmManager,
		plansRepository:  plansRepository,
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
