package rocket_league

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func handleGet(services types.Services, channelId string) (*model.ChannelsIntegrationsData, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetChannelIntegration(channelId, model.IntegrationServiceRocketLeague, services.DB)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	if integration == nil {
		return nil, nil
	}

	return integration.Data, nil
}

func handlePost(services types.Services, channelId string, dto *createOrUpdateDTO) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetChannelIntegration(channelId, model.IntegrationServiceRocketLeague, services.DB)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	newData := &model.ChannelsIntegrationsData{
		UserName: &dto.Username,
		Code:     &dto.Code,
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.DB.
			Where("service = ?", model.IntegrationServiceRocketLeague).
			First(&neededIntegration).
			Error
		if err != nil {
			logger.Error(err)
			return fiber.NewError(
				http.StatusInternalServerError,
				"seems like Rocket League not enabled on our side",
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
			Data:          newData,
			Integration:   nil,
		}
	}

	integration.Data = newData
	err = services.DB.Save(integration).Error
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}
