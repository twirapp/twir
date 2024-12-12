package timers

import (
	"errors"

	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
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
}

func New(opts Opts) *Service {
	return &Service{
		gorm:             opts.Gorm,
		logger:           opts.Logger,
		twirbus:          opts.TwirBus,
		timersrepository: opts.TimersRepository,
	}
}

type Service struct {
	gorm             *gorm.DB
	logger           logger.Logger
	twirbus          *buscore.Bus
	timersrepository timersrepository.Repository
}

const MaxTimersPerChannel = 10

var ErrTimerNotFound = errors.New("timer not found")

func (c *Service) dbToModel(m timersmodel.Timer) entity.Timer {
	responses := make([]entity.Response, 0, len(m.Responses))
	for _, r := range m.Responses {
		responses = append(
			responses,
			entity.Response{
				ID:         r.ID,
				Text:       r.Text,
				IsAnnounce: r.IsAnnounce,
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
