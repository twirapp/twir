package streamlabs

import (
	"context"
	"fmt"
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
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGetAuth(services types.Services) (*string, error) {
	integration := model.Integrations{}
	err := services.DB.Where(`"service" = ?`, "STREAMLABS").First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(
			404,
			"streamlabs not enabled on our side. Please be patient.",
		)
	}

	url, _ := url.Parse("https://www.streamlabs.com/api/v1.0/authorize")

	q := url.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", "socket.token donations.read")
	q.Add("redirect_uri", integration.RedirectURL.String)
	url.RawQuery = q.Encode()

	str := url.String()

	return &str, nil
}

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrationsData, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	integration, err := helpers.GetIntegration(channelId, "STREAMLABS", services.DB)
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
	StreamLabs struct {
		ID          int    `json:"id"`
		DisplayName string `json:"display_name"`
		ThumbNail   string `json:"thumbnail"`
	} `json:"streamlabs"`
}

func handlePost(channelId string, dto *tokenDto, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	channelIntegration, err := helpers.GetIntegration(channelId, "STREAMLABS", services.DB)
	if err != nil {
		logger.Error(err)
		return err
	}

	neededIntegration := model.Integrations{}
	err = services.DB.
		Where("service = ?", "STREAMLABS").
		First(&neededIntegration).
		Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(
			http.StatusInternalServerError,
			"seems like streamlabs not enabled on our side",
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
		Post("https://streamlabs.com/api/v1.0/token")
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
		Get(fmt.Sprintf("https://streamlabs.com/api/v1.0/user?access_token=%s", data.AccessToken))

	if err != nil || !profileResp.IsSuccess() {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot get profile")
	}

	channelIntegration.AccessToken = null.StringFrom(data.AccessToken)
	channelIntegration.RefreshToken = null.StringFrom(data.RefreshToken)
	channelIntegration.Data = &model.ChannelsIntegrationsData{
		UserId: lo.ToPtr(fmt.Sprint(profile.StreamLabs.ID)),
		Name:   lo.ToPtr(profile.StreamLabs.DisplayName),
		Avatar: lo.ToPtr(profile.StreamLabs.ThumbNail),
	}

	err = services.DB.
		Save(channelIntegration).Error

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update integration")
	}

	sendGrpcEvent(channelIntegration.ID, channelIntegration.Enabled, services)

	return nil
}

func sendGrpcEvent(integrationId string, isAdd bool, services types.Services) {
	if isAdd {
		services.IntegrationsGrpc.AddIntegration(context.Background(), &integrations.Request{
			Id: integrationId,
		})
	} else {
		services.IntegrationsGrpc.RemoveIntegration(context.Background(), &integrations.Request{
			Id: integrationId,
		})
	}
}

func handleLogout(channelId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	integration, err := helpers.GetIntegration(channelId, "STREAMLABS", services.DB)
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
