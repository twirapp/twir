package giveaways

import (
	"github.com/satont/twir/apps/api/internal/types"
	model "github.com/satont/twir/libs/gomodels"
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

func handlePost(channelId string, dto *postGiveawayDto, services types.Services) (*model.ChannelGiveaway, error) {
	// logger := do.MustInvoke[interfaces.Logger](di.Provider)

	// giveaway := &model.ChannelGiveaway{
	// 	ID: uuid.NewV4().String(),
	// }
	return nil, nil
}
