package donationalerts

import (
	"context"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"net/http"
	"net/url"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
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

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrationsData, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	integration, err := helpers.GetIntegration(channelId, "DONATIONALERTS", services.DB)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	if integration == nil {
		return nil, nil
	}

	return integration.Data, nil
}

type tokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type profileResponse struct {
	Data struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Avatar string `json:"avatar"`
	} `json:"data"`
}

func handlePost(channelId string, dto *tokenDto, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	channelIntegration, err := helpers.GetIntegration(channelId, "DONATIONALERTS", services.DB)
	if err != nil {
		logger.Error(err)
		return err
	}

	neededIntegration := model.Integrations{}
	err = services.DB.
		Where("service = ?", "DONATIONALERTS").
		First(&neededIntegration).
		Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(
			http.StatusInternalServerError,
			"seems like donationalerts not enabled on our side",
		)
	}

	if channelIntegration == nil {
		channelIntegration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
			Enabled:       true,
		}
	}

	data := tokensResponse{}
	resp, err := req.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     neededIntegration.ClientID.String,
			"client_secret": neededIntegration.ClientSecret.String,
			"redirect_uri":  neededIntegration.RedirectURL.String,
			"code":          dto.Code,
		}).
		SetResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://www.donationalerts.com/oauth/token")
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get tokens")
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
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get profile")
	}

	if channelIntegration == nil {
		channelIntegration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
			Enabled:       true,
		}
	}

	channelIntegration.Data = &model.ChannelsIntegrationsData{
		Name:   &profile.Data.Name,
		Code:   &profile.Data.Code,
		Avatar: &profile.Data.Avatar,
	}

	channelIntegration.AccessToken = null.StringFrom(data.AccessToken)
	channelIntegration.RefreshToken = null.StringFrom(data.RefreshToken)

	err = services.DB.Save(channelIntegration).Error

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update integration")
	}

	sendGrpcEvent(channelIntegration.ID, channelIntegration.Enabled)

	return nil
}

func sendGrpcEvent(integrationId string, isAdd bool) {
	grpc := do.MustInvoke[integrations.IntegrationsClient](di.Injector)
	if isAdd {
		grpc.AddIntegration(context.Background(), &integrations.Request{
			Id: integrationId,
		})
	} else {
		grpc.RemoveIntegration(context.Background(), &integrations.Request{
			Id: integrationId,
		})
	}
}

func handleLogout(channelId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	integration, err := helpers.GetIntegration(channelId, "DONATIONALERTS", services.DB)
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
