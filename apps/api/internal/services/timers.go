package services

import (
	"errors"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var notFoundError = errors.New("timer not found")

type TimersService struct {
	db     *gorm.DB
	logger interfaces.Logger
}

func NewTimersService() *TimersService {
	service := &TimersService{
		db:     do.MustInvoke[*gorm.DB](di.Injector),
		logger: do.MustInvoke[interfaces.Logger](di.Injector),
	}

	return service
}

func (c *TimersService) FindOneById(id string) (*model.ChannelsTimers, error) {
	timer := model.ChannelsTimers{}
	err := c.db.Where("id = ?", id).Preload("Responses").Find(&timer).Error

	if err != nil {
		c.logger.Error(err)
		return nil, notFoundError
	}

	if timer.ID == "" {
		return nil, notFoundError
	}

	return &timer, nil
}

func (c *TimersService) FindManyByChannelId(channelId string) ([]model.ChannelsTimers, error) {
	var timers []model.ChannelsTimers
	err := c.db.Where(`"channelId" = ?`, channelId).Preload("Responses").Find(&timers).Error

	if err != nil {
		c.logger.Error(err)
		return nil, notFoundError
	}

	return timers, nil
}

func (c *TimersService) Create(
	timer model.ChannelsTimers,
	responses []model.ChannelsTimersResponses,
) (*model.ChannelsTimers, error) {
	if len(responses) == 0 {
		return nil, errors.New("responses cannot be empty")
	}
	timer.ID = uuid.NewV4().String()

	err := c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&timer).Error; err != nil {
			return err
		}

		for _, response := range responses {
			response.TimerID = timer.ID
			response.ID = uuid.NewV4().String()
			if err := tx.Save(&response).Error; err != nil {
				return err
			}
			timer.Responses = append(timer.Responses, response)
		}

		return nil
	})

	if err != nil {
		c.logger.Error(err)
		return nil, errors.New("internal error")
	}

	return &timer, nil
}

func (c *TimersService) Delete(id string) error {
	timer, err := c.FindOneById(id)
	if err != nil {
		return err
	}
	err = c.db.Delete(&timer).Error
	if err != nil {
		c.logger.Error(err)
		return errors.New("cannot delete timer")
	}
	return nil
}

func (c *TimersService) Update(
	timerId string,
	data model.ChannelsTimers,
	responses []model.ChannelsTimersResponses,
) (*model.ChannelsTimers, error) {
	if len(responses) == 0 {
		return nil, errors.New("responses cannot be empty")
	}

	timer, err := c.FindOneById(timerId)
	if err != nil {
		return nil, err
	}

	timer.Enabled = data.Enabled
	timer.LastTriggerMessageNumber = 0
	timer.MessageInterval = data.MessageInterval
	timer.Name = data.Name
	timer.TimeInterval = data.TimeInterval

	timer.Responses = []model.ChannelsTimersResponses{}

	err = c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&timer).Error; err != nil {
			return err
		}

		err := tx.Where(`"timerId" = ?`, timer.ID).Delete(&model.ChannelsTimersResponses{}).Error
		if err != nil {
			return err
		}

		for _, response := range responses {
			response.TimerID = timer.ID
			response.ID = uuid.NewV4().String()
			if err := tx.Save(&response).Error; err != nil {
				return err
			}
			timer.Responses = append(timer.Responses, response)
		}

		return nil
	})

	if err != nil {
		c.logger.Error(err)
		return nil, errors.New("internal error")
	}

	return timer, nil
}

func (c *TimersService) SetEnabled(timerId string, enabled bool) (*model.ChannelsTimers, error) {
	timer, err := c.FindOneById(timerId)
	if err != nil {
		return nil, err
	}

	timer.Enabled = enabled
	if err := c.db.Save(&timer).Error; err != nil {
		c.logger.Error(err)
		return nil, errors.New("internal error")
	}

	return timer, nil
}
