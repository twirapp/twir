package donate_stream

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

func handleGet(services types.Services, channelId string) (string, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetChannelIntegration(channelId, model.IntegrationServiceDonateStream, services.DB)
	if err != nil {
		logger.Error(err)
		return "", nil
	}

	if integration == nil {
		neededIntegration := &model.Integrations{}
		err = services.DB.
			Where("service = ?", model.IntegrationServiceDonateStream).
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

func handlePost(services types.Services, integrationId, secret string) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	redis := do.MustInvoke[*redis.Client](di.Provider)

	err := redis.Set(context.Background(), "donate_stream_confirmation"+integrationId, secret, 1*time.Hour).Err()
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
