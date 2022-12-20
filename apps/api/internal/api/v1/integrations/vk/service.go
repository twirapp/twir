package vk

import (
	"net/http"
	"net/url"

	"github.com/davecgh/go-spew/spew"
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	req "github.com/imroc/req/v3"
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
	integration, err := helpers.GetIntegration(channelId, "VK", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, err
	}

	data := profileResponse{}
	_, err = req.R().
		SetQueryParams(map[string]string{
			"v":            "5.131",
			"fields":       "photo_max_orig",
			"access_token": integration.AccessToken.String,
		}).
		SetResult(&data).
		Get("https://api.vk.com/method/users.get")
	spew.Dump(data)

	if err != nil {
		services.Logger.Sugar().Error(err)
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
			Integration:   &neededIntegration,
			Enabled:       true,
		}
	}

	data := tokensResponse{}
	_, err = req.R().
		SetQueryParams(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     integration.Integration.ClientID.String,
			"client_secret": integration.Integration.ClientSecret.String,
			"redirect_uri":  integration.Integration.RedirectURL.String,
			"code":          dto.Code,
		}).
		SetResult(&data).
		Get("https://oauth.vk.com/access_token")

	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	integration.Enabled = true
	integration.AccessToken = null.StringFrom(data.AccessToken)

	if err = services.DB.Save(integration).Error; err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update faceit data")
	}

	return nil
}
