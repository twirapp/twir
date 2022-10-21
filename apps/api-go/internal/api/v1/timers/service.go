package timers

import (
	model "tsuwari/models"

	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func handleGet(channelId string, services types.Services) []model.ChannelsTimers {
	timers := []model.ChannelsTimers{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Find(&timers).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil
	}
	return timers
}
