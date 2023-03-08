package donatepay

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	uuid "github.com/satori/go.uuid"
)

const integrationName = "DONATEPAY"

func handleGet(services *types.Services, channelId string) (*string, error) {
	

	integration, err := helpers.GetIntegration(channelId, integrationName, services.Gorm)
	if err != nil {
		services.Logger.Error(err)
		return nil, nil
	}

	if integration == nil {
		return nil, nil
	}

	return &integration.APIKey.String, nil
}

func handlePost(services *types.Services, channelId string, dto *createOrUpdateDTO) error {

	integration, err := helpers.GetIntegration(channelId, integrationName, services.Gorm)
	if err != nil {
		services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.Gorm.
			Where("service = ?", integrationName).
			First(&neededIntegration).
			Error
		if err != nil {
			services.Logger.Error(err)
			return fiber.NewError(
				http.StatusInternalServerError,
				"seems like donatepay not enabled on our side",
			)
		}

		integration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			Enabled:       true,
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
			AccessToken:   null.String{},
			RefreshToken:  null.String{},
			ClientID:      null.String{},
			ClientSecret:  null.String{},
			Data:          nil,
			Integration:   nil,
		}
	}

	integration.APIKey = null.StringFrom(dto.ApiKey)
	err = services.Gorm.Save(integration).Error
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if len(integration.APIKey.String) > 0 {
		services.Grpc.Integrations.AddIntegration(context.Background(), &integrations.Request{
			Id: integration.ID,
		})
	} else {
		services.Grpc.Integrations.RemoveIntegration(context.Background(), &integrations.Request{
			Id: integration.ID,
		})
	}

	return nil
}
