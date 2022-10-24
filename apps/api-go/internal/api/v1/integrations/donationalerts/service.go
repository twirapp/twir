package donationalerts

import (
	"net/url"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	"github.com/satont/tsuwari/libs/nats/integrations"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func handleGetAuth(services types.Services) (*string, error) {
	integration := model.Integrations{}
	err := services.DB.Where(`"service" = ?`, "DONATIONALERTS").First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(
			404,
			"donationalerts not enabled on our side. Please be patient.",
		)
	}

	url, _ := url.Parse("https://www.donationalerts.com/oauth/authorize")

	q := url.Query()
	q.Add("client_id", integration.ClientID.String)
	q.Add("response_type", "code")
	q.Add("scope", "oauth-user-show oauth-donation-subscribe")
	q.Add("redirect_uri", integration.RedirectURL.String)
	url.RawQuery = q.Encode()

	str := url.String()

	return &str, nil
}

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrations, error) {
	integration := model.ChannelsIntegrations{}
	err := services.DB.Where(`"channelId" = ?`, channelId).
		Preload("Integration").
		First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "internal error")
	}
	return &integration, nil
}

func handlePatch(channelId string, dto *donationAlertsDto, services types.Services) error {
	integration := model.ChannelsIntegrations{}
	err := services.DB.Where(`"channelId" = ?`, channelId).
		First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(404, "integration not found")
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "internal error")
	}

	integration.Enabled = *dto.Enabled
	services.DB.Save(&integration)

	bytes := []byte{}
	if *dto.Enabled {
		bytes, _ = proto.Marshal(&integrations.AddIntegration{Id: integration.ID})
	} else {
		bytes, _ = proto.Marshal(&integrations.RemoveIntegration{Id: integration.ID})
	}
	services.Nats.Publish(
		lo.If(*dto.Enabled, integrations.SUBJECTS_ADD_INTEGRATION).
			Else(integrations.SUBJECTS_REMOVE_INTEGRATION),
		bytes,
	)

	return nil
}
