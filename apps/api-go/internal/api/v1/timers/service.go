package timers

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	uuid "github.com/satori/go.uuid"
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
