package faceit

import (
	"net/http"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrations, error) {
	integration, err := helpers.GetIntegration(channelId, "FACEIT", services.DB)
	if err != nil {
		return nil, err
	}

	return integration, nil
}

func handlePost(channelId string, dto *faceitUpdateDto, services types.Services) error {
	integration, err := helpers.GetIntegration(channelId, "FACEIT", services.DB)
	if err != nil {
		return err
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.DB.
			Where("service = ?", "FACEIT").
			First(&neededIntegration).
			Error
		if err != nil {
			return fiber.NewError(
				http.StatusInternalServerError,
				"seems like faceit not enabled on our side",
			)
		}

		integration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
		}
	}

	if dto.Data.Game == nil {
		dto.Data.Game = lo.ToPtr("csgo")
	}

	dto.Data.Game = lo.ToPtr(strings.ToLower(*dto.Data.Game))

	integration.Enabled = *dto.Enabled
	integration.Data = &model.ChannelsIntegrationsData{
		Game:     dto.Data.Game,
		UserName: &dto.Data.UserName,
	}

	if err = services.DB.Save(integration).Error; err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update faceit data")
	}

	return nil
}
