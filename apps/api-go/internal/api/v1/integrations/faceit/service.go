package faceit

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrations, error) {
	integration := model.ChannelsIntegrations{}

	err := services.DB.
		Preload("Integration").
		Joins(`JOIN integrations i on i.id = channels_integrations."integrationId"`).
		Where(`"channels_integrations"."channelId" = ? AND i.service = ?`, channelId, "FACEIT").
		First(&integration).
		Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(404, "integration not found")
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "internal error")
	}

	return &integration, nil
}
