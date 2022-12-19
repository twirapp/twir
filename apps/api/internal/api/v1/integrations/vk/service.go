package vk

import (
	"net/http"
	"net/url"

	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

type VK struct {
	model.ChannelsIntegrations
	Data map[string]any `json:"data"`
}

func handleGetAuth(services types.Services) (*string, error) {
	integration := model.Integrations{}
	err := services.DB.Where(`"service" = ?`, "VK").First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(
			404,
			"donationalerts not enabled on our side. Please be patient.",
		)
	}

	url, _ := url.Parse("https://oauth.vk.com/authorize")

	q := url.Query()
	q.Add("client_id", integration.ClientID.String)
	q.Add("display", "page")
	q.Add("response_type", "code")
	q.Add("scope", "status offline")
	q.Add("redirect_uri", integration.RedirectURL.String)
	url.RawQuery = q.Encode()

	str := url.String()

	return &str, nil
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
