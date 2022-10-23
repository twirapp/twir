package timers

import (
	"errors"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) []model.ChannelsTimers {
	timers := []model.ChannelsTimers{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Preload("Responses").Find(&timers).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil
	}
	return timers
}

func handlePost(
	channelId string,
	dto *timerDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timer := model.ChannelsTimers{
		ID:                       uuid.NewV4().String(),
		ChannelID:                channelId,
		Name:                     dto.Name,
		Enabled:                  *dto.Enabled,
		TimeInterval:             int32(dto.TimeInterval),
		MessageInterval:          int32(dto.MessageInterval),
		LastTriggerMessageNumber: 0,
	}

	err := services.DB.Save(&timer).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot create timer")
	}

	timerResponses := []model.ChannelsTimersResponses{}

	for _, t := range dto.Responses {
		response := model.ChannelsTimersResponses{
			ID:         uuid.NewV4().String(),
			Text:       t.Text,
			IsAnnounce: *t.IsAnnounce,
			TimerID:    timer.ID,
		}

		err := services.DB.Save(&response).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			services.DB.Where(`"id" = ?`, timer.ID).Delete(&model.ChannelsTimers{})

			return nil, fiber.NewError(500, "cannot create timer responses")
		}

		timerResponses = append(timerResponses, response)
	}

	timer.Responses = &timerResponses
	return &timer, nil
}

func handleDelete(timerId string, services types.Services) error {
	timer := model.ChannelsTimers{}
	err := services.DB.Where("id = ?", timerId).First(&timer).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return fiber.NewError(404, "timer not found")
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "timer not found")
	}

	err = services.DB.Delete(&timer).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot delete timer")
	}

	return nil
}

func handlePut(
	channelId string,
	timerId string,
	dto *timerDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timer := model.ChannelsTimers{}

	err := services.DB.Where(`"channelId" = ? AND "id" = ?`, channelId, timerId).
		Preload("Responses").
		First(&timer).
		Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(401, "timer with that id not found")
	}

	if err != nil {
		return nil, fiber.NewError(500, "internal error")
	}

	timer.Enabled = *dto.Enabled
	timer.LastTriggerMessageNumber = 0
	timer.MessageInterval = int32(dto.MessageInterval)
	timer.Name = dto.Name
	timer.TimeInterval = int32(dto.TimeInterval)

	err = services.DB.Select("*").Updates(&timer).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot update timer")
	}

	err = services.DB.Transaction(func(tx *gorm.DB) error {
		for _, response := range *timer.Responses {
			err = tx.Where("id = ?", response.ID).Delete(&model.ChannelsTimersResponses{}).Error
			if err != nil {
				services.Logger.Sugar().Error(err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "internal error when deleting responses")
	}

	newResponses := []model.ChannelsTimersResponses{}
	for _, response := range dto.Responses {
		r := model.ChannelsTimersResponses{
			ID:         uuid.NewV4().String(),
			Text:       response.Text,
			IsAnnounce: *response.IsAnnounce,
			TimerID:    timer.ID,
		}
		services.DB.Save(&r)
		newResponses = append(newResponses, r)
	}

	timer.Responses = &newResponses

	return &timer, nil
}
