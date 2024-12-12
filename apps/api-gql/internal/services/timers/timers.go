package timers

import (
	"github.com/google/uuid"
	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers/model"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Gorm    *gorm.DB
	Logger  logger.Logger
	TwirBus *buscore.Bus
}

func New(opts Opts) *Service {
	return &Service{
		gorm:    opts.Gorm,
		logger:  opts.Logger,
		twirbus: opts.TwirBus,
	}
}

type Service struct {
	gorm    *gorm.DB
	logger  logger.Logger
	twirbus *buscore.Bus
}

const MaxTimersPerChannel = 10

func (c *Service) dbToModel(m dbmodels.ChannelsTimers) model.Timer {
	responses := make([]model.Response, 0, len(m.Responses))
	for _, r := range m.Responses {
		responses = append(
			responses,
			model.Response{
				ID:         uuid.MustParse(r.ID),
				Text:       r.Text,
				IsAnnounce: r.IsAnnounce,
			},
		)
	}

	return model.Timer{
		ID:                       uuid.MustParse(m.ID),
		ChannelID:                m.ChannelID,
		Name:                     m.Name,
		Enabled:                  m.Enabled,
		TimeInterval:             int(m.TimeInterval),
		MessageInterval:          int(m.MessageInterval),
		LastTriggerMessageNumber: int(m.LastTriggerMessageNumber),
		Responses:                responses,
	}
}
