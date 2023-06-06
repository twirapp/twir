package giveaways

import (
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"go.uber.org/zap"
)

func handleGetAll(channelId string, services types.Services) []*model.ChannelGiveaway {
	var giveaways []*model.ChannelGiveaway

	err := services.DB.
		Where(`channel_id" = ?`, channelId).
		Order(`created_at desc`).
		Find(&giveaways).Error
	if err != nil {
		zap.S().Error(err)
		return nil
	}

	return giveaways
}
