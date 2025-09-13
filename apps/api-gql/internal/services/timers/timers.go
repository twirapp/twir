package timers

import (
	"errors"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/logger"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	timersmodel "github.com/twirapp/twir/libs/repositories/timers/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm             *gorm.DB
	Logger           logger.Logger
	TwirBus          *buscore.Bus
	TimersRepository timersrepository.Repository
	TrmManager       trm.Manager
}

func New(opts Opts) *Service {
	return &Service{
		gorm:             opts.Gorm,
		logger:           opts.Logger,
		twirbus:          opts.TwirBus,
		timersRepository: opts.TimersRepository,
		trmManager:       opts.TrmManager,
	}
}

type Service struct {
	gorm             *gorm.DB
	logger           logger.Logger
	twirbus          *buscore.Bus
	timersRepository timersrepository.Repository
	trmManager       trm.Manager
}

const MaxPerChannel = 10

var ErrTimerNotFound = errors.New("timer not found")

func (c *Service) dbToModel(m timersmodel.Timer) entity.Timer {
	responses := make([]entity.Response, 0, len(m.Responses))
	for _, r := range m.Responses {
		responses = append(
			responses,
			entity.Response{
				ID:            r.ID,
				Text:          r.Text,
				IsAnnounce:    r.IsAnnounce,
				Count:         r.Count,
				AnnounceColor: bots.AnnounceColor(r.AnnounceColor),
			},
		)
	}

	return entity.Timer{
		ID:                       m.ID,
		ChannelID:                m.ChannelID,
		Name:                     m.Name,
		Enabled:                  m.Enabled,
		TimeInterval:             m.TimeInterval,
		MessageInterval:          m.MessageInterval,
		LastTriggerMessageNumber: m.LastTriggerMessageNumber,
		Responses:                responses,
	}
}
