package donatello

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

const integrationName = "DONATELLO"

func handleGet(services types.Services, channelId string) (*string, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	integration, err := helpers.GetIntegration(channelId, integrationName, services.DB)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	if integration == nil {
		return nil, nil
	}

	return &integration.APIKey.String, nil
}

func handlePost(services types.Services, channelId string, dto *createOrUpdateDTO) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)
	integrationsGrpc := do.MustInvoke[integrations.IntegrationsClient](di.Injector)

	integration, err := helpers.GetIntegration(channelId, integrationName, services.DB)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.DB.
			Where("service = ?", integrationName).
			First(&neededIntegration).
			Error
		if err != nil {
			logger.Error(err)
			return fiber.NewError(
				http.StatusInternalServerError,
				"seems like donatello not enabled on our side",
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
	err = services.DB.Save(integration).Error
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if len(integration.APIKey.String) > 0 {
		integrationsGrpc.AddIntegration(context.Background(), &integrations.Request{
			Id: integration.ID,
		})
	} else {
		integrationsGrpc.RemoveIntegration(context.Background(), &integrations.Request{
			Id: integration.ID,
		})
	}

	return nil
}
