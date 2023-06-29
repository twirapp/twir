package vk

import (
	"net/http"
	"net/url"

	"github.com/samber/do"
	"github.com/satont/twir/apps/api/internal/di"
	"github.com/satont/twir/apps/api/internal/interfaces"

	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
	"github.com/satont/twir/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/twir/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

type VK struct {
	model.ChannelsIntegrations
	Data map[string]any `json:"data"`
}

func handleGetAuth(services types.Services) (*string, error) {
	integration := model.Integrations{}
	err := services.DB.Where(`"service" = ?`, model.IntegrationServiceVK).First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(
			404,
			"vk not enabled on our side. Please be patient.",
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

type profile struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhotoMaxOrig string `json:"photo_max_orig"`
}

type profileResponse struct {
	Response []profile `json:"response"`
	Error    *struct {
		Code int    `json:"error_code"`
		Msg  string `json:"error_msg"`
	}
}

func handleGet(channelId string, services types.Services) (*profile, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetChannelIntegration(channelId, model.IntegrationServiceVK, services.DB)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if integration == nil {
		return nil, nil
	}

	data := profileResponse{}
	_, err = req.R().
		SetQueryParams(
			map[string]string{
				"v":            "5.131",
				"fields":       "photo_max_orig",
				"access_token": integration.AccessToken.String,
			},
		).
		SetResult(&data).
		Get("https://api.vk.com/method/users.get")

	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if data.Error != nil || len(data.Response) == 0 {
		return nil, nil
	}

	return &data.Response[0], nil
}

type tokensResponse struct {
	AccessToken string `json:"access_token"`
}

func handlePost(channelId string, dto *vkDto, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetChannelIntegration(channelId, model.IntegrationServiceVK, services.DB)
	if err != nil {
		logger.Error(err)
		return err
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.DB.
			Where("service = ?", model.IntegrationServiceVK).
			First(&neededIntegration).
			Error
		if err != nil {
			logger.Error(err)
			return fiber.NewError(
				http.StatusInternalServerError,
				"seems like VK not enabled on our side",
			)
		}

		integration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
			Integration:   &neededIntegration,
			Enabled:       true,
		}
	}

	data := tokensResponse{}
	_, err = req.R().
		SetQueryParams(
			map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     integration.Integration.ClientID.String,
				"client_secret": integration.Integration.ClientSecret.String,
				"redirect_uri":  integration.Integration.RedirectURL.String,
				"code":          dto.Code,
			},
		).
		SetResult(&data).
		Get("https://oauth.vk.com/access_token")

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	integration.Enabled = true
	integration.AccessToken = null.StringFrom(data.AccessToken)

	if err = services.DB.Save(integration).Error; err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update faceit data")
	}

	return nil
}

func handleLogout(channelId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetChannelIntegration(channelId, model.IntegrationServiceVK, services.DB)
	if err != nil {
		logger.Error(err)
		return err
	}
	if integration == nil {
		return fiber.NewError(http.StatusNotFound, "integration not found")
	}

	err = services.DB.Delete(&integration).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
