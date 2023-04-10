package faceit

import (
	"encoding/base64"
	"fmt"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrationsData, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetIntegration(channelId, "FACEIT", services.DB)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	if integration == nil {
		return nil, nil
	}

	return integration.Data, nil
}

func handleGetAuth(services types.Services) (*string, error) {
	integration := model.Integrations{}
	err := services.DB.Where(`"service" = ?`, "FACEIT").First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(
			404,
			"spotify not enabled on our side. Please be patient.",
		)
	}

	url, _ := url.Parse("https://cdn.faceit.com/widgets/sso/index.html")

	q := url.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("redirect_popup", integration.RedirectURL.String)
	url.RawQuery = q.Encode()

	str := url.String()

	return &str, nil
}

func handlePost(channelId string, dto *tokenDto, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	channelIntegration, err := helpers.GetIntegration(channelId, "FACEIT", services.DB)

	if err != nil {
		logger.Error(err)
		return err
	}

	neededIntegration := model.Integrations{}
	err = services.DB.
		Where("service = ?", "FACEIT").
		First(&neededIntegration).
		Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(
			http.StatusInternalServerError,
			"seems like faceit not enabled on our side",
		)
	}

	if channelIntegration == nil {
		channelIntegration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
			Enabled:       true,
			Integration:   &neededIntegration,
		}
	}

	tokensData := make(map[string]any)
	token := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf(
			"%s:%s",
			neededIntegration.ClientID.String,
			neededIntegration.ClientSecret.String,
		),
	))
	resp, err := req.R().
		SetFormData(map[string]string{
			"grant_type": "authorization_code",
			"code":       dto.Code,
		}).
		SetSuccessResult(&tokensData).
		SetHeader("Authorization", fmt.Sprintf("Basic %s", token)).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://api.faceit.com/auth/v1/oauth/token")
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "cannot get tokens")
	}

	if !resp.IsSuccessState() {
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))
		fmt.Println(resp.StatusCode)
		return fiber.NewError(401, "seems like code is invalid")
	}

	channelIntegration.AccessToken = null.StringFrom(tokensData["access_token"].(string))
	channelIntegration.RefreshToken = null.StringFrom(tokensData["refresh_token"].(string))

	userInfoResult := make(map[string]any)
	resp, err = req.R().
		SetBearerAuthToken(channelIntegration.AccessToken.String).
		SetResult(&userInfoResult).
		Get("https://api.faceit.com/auth/v1/resources/userinfo")

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "cannot get tokens")
	}

	if !resp.IsSuccess() {
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))
		return fiber.NewError(401, "seems response of profile is incorrect")
	}

	integrationData := model.ChannelsIntegrationsData{
		UserId: lo.ToPtr(userInfoResult["guid"].(string)),
		Name:   lo.ToPtr(userInfoResult["nickname"].(string)),
		Game:   lo.ToPtr("csgo"),
	}

	profileResult := make(map[string]any)
	resp, err = req.R().
		SetBearerAuthToken(channelIntegration.Integration.APIKey.String).
		SetResult(&profileResult).
		Get("https://open.faceit.com/data/v4/players/" + *integrationData.UserId)

	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "cannot get tokens")
	}

	if !resp.IsSuccess() {
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))
		return fiber.NewError(401, "seems response of profile is incorrect")
	}

	integrationData.Avatar = lo.ToPtr(profileResult["avatar"].(string))

	channelIntegration.Data = &integrationData

	if err = services.DB.Save(channelIntegration).Error; err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update integration")
	}

	return nil
}

func handleLogout(channelId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	integration, err := helpers.GetIntegration(channelId, "FACEIT", services.DB)
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
