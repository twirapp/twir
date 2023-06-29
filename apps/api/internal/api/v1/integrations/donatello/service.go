package donatello

import (
	"github.com/guregu/null"
	"github.com/samber/do"
	"github.com/satont/twir/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/twir/apps/api/internal/di"
	"github.com/satont/twir/apps/api/internal/interfaces"
	"github.com/satont/twir/apps/api/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func handleGet(services types.Services, channelId string) (string, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetChannelIntegration(channelId, model.IntegrationServiceDonatello, services.DB)
	if err != nil {
		logger.Error(err)
		return "", nil
	}

	if integration == nil {
		neededIntegration := &model.Integrations{}
		err = services.DB.
			Where("service = ?", model.IntegrationServiceDonatello).
			First(neededIntegration).
			Error

		if err != nil {
			logger.Error(err)
			return "", err
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
		}

		err = services.DB.Save(integration).Error
		if err != nil {
			logger.Error(err)
			return "", err
		}

		return integration.ID, nil
	}

	return integration.ID, nil
}
