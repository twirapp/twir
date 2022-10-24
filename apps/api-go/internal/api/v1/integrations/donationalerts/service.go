package donationalerts

import (
	"net/url"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	req "github.com/imroc/req/v3"
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
	err := services.DB.
		Preload("Integration").
		Joins(`JOIN integrations i on i.id = channels_integrations."integrationId"`).
		Where(`"channels_integrations"."channelId" = ? AND i.service = ?`, channelId, "DONATIONALERTS").
		First(&integration).
		Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "internal error")
	}
	return &integration, nil
}

func handlePatch(
	channelId string,
	dto *donationAlertsDto,
	services types.Services,
) (*model.ChannelsIntegrations, error) {
	integration := model.ChannelsIntegrations{}
	err := services.DB.Where(`"channelId" = ?`, channelId).
		First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(404, "integration not found")
	}
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "internal error")
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

	return &integration, nil
}

type tokensResponse struct {
	AccessToken   string `json:"access_token"`
	RerfreshToken string `json:"rerfresh_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
}

type profileResponse struct {
	Data struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Avatar string `json:"avatar"`
	} `json:"data"`
}

func handlePost(channelId string, dto *tokenDto, services types.Services) error {
	integration := model.Integrations{}
	err := services.DB.Where(`"service" = ?`, "DONATIONALERTS").First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return fiber.NewError(
			404,
			"donationalerts not enabled on our side. Please be patient.",
		)
	}

	data := tokensResponse{}
	resp, err := req.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     integration.ClientID.String,
			"client_secret": integration.ClientSecret.String,
			"redirect_uri":  integration.RedirectURL.String,
			"code":          dto.Code,
		}).
		SetResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://www.donationalerts.com/oauth/token")
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot get tokens")
	}
	if !resp.IsSuccess() {
		return fiber.NewError(401, "seems like code is invalid")
	}

	profile := profileResponse{}
	profileResp, err := req.R().
		SetResult(&profile).
		SetBearerAuthToken(data.AccessToken).
		Get("https://www.donationalerts.com/api/v1/user/oauth")

	if err != nil || !profileResp.IsSuccess() {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot get profile")
	}

	err = services.DB.
		Model(&model.ChannelsIntegrations{}).
		Where(`"integrationId" = ?`, integration.ID).
		Updates(map[string]interface{}{
			"accessToken":  data.AccessToken,
			"refreshToken": data.RerfreshToken,
			"data": map[string]string{
				"name":   profile.Data.Name,
				"code":   profile.Data.Code,
				"avatar": profile.Data.Avatar,
			},
		}).Error

	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot update integration")
	}

	return nil
}
