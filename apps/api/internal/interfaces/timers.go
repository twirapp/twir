package interfaces

import model "github.com/satont/twir/libs/gomodels"

type TimersService interface {
	FindOneById(id string) (*model.ChannelsTimers, error)
	FindManyByChannelId(channelId string) ([]model.ChannelsTimers, error)
	Create(data model.ChannelsTimers, responses []model.ChannelsTimersResponses) (*model.ChannelsTimers, error)
	Delete(id string) error
	Update(timerId string, data model.ChannelsTimers, responses []model.ChannelsTimersResponses) (
		*model.ChannelsTimers, error,
	)
	SetEnabled(timerId string, enabled bool) (*model.ChannelsTimers, error)
}
