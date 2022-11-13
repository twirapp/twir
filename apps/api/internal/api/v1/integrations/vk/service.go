package vk

import (
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

type VK struct {
	model.ChannelsIntegrations
	Data map[string]any `json:"data"`
}

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrations, error) {
	integration, err := helpers.GetIntegration(channelId, "VK", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, err
	}

	return integration, nil
}

func handlePost(channelId string, dto *vkDto, services types.Services) error {
	integration, err := helpers.GetIntegration(channelId, "VK", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return err
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.DB.
			Where("service = ?", "VK").
			First(&neededIntegration).
			Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			return fiber.NewError(
				http.StatusInternalServerError,
				"seems like VK not enabled on our side",
			)
		}

		integration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
		}
	}

	integration.Enabled = *dto.Enabled
	integration.Data = &model.ChannelsIntegrationsData{
		UserId: &dto.Data.UserId,
	}

	if err = services.DB.Save(integration).Error; err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update faceit data")
	}

	return nil
}
