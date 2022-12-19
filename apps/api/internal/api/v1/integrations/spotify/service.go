package spotify

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	req "github.com/imroc/req/v3"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	spotify "github.com/satont/tsuwari/libs/integrations/spotify"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGetAuth(services types.Services) (*string, error) {
	integration := model.Integrations{}
	err := services.DB.Where(`"service" = ?`, "SPOTIFY").First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(
			404,
			"spotify not enabled on our side. Please be patient.",
		)
	}

	url, _ := url.Parse("https://accounts.spotify.com/authorize")

	q := url.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", "user-read-currently-playing")
	q.Add("redirect_uri", integration.RedirectURL.String)
	url.RawQuery = q.Encode()

	str := url.String()

	return &str, nil
}

func handleGet(channelId string, services types.Services) (*model.ChannelsIntegrations, error) {
	integration, err := helpers.GetIntegration(channelId, "SPOTIFY", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, err
	}

	return integration, nil
}

func handlePatch(
	channelId string,
	dto *spotifyDto,
	services types.Services,
) (*model.ChannelsIntegrations, error) {
	integration, err := helpers.GetIntegration(channelId, "SPOTIFY", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, err
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.DB.
			Where("service = ?", "SPOTIFY").
			First(&neededIntegration).
			Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"seems like spotify not enabled on our side",
			)
		}

		integration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
		}
	}

	integration.Enabled = *dto.Enabled
	services.DB.Save(&integration)

	return integration, nil
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
	channelIntegration, err := helpers.GetIntegration(channelId, "SPOTIFY", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return err
	}

	neededIntegration := model.Integrations{}
	err = services.DB.
		Where("service = ?", "SPOTIFY").
		First(&neededIntegration).
		Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(
			http.StatusInternalServerError,
			"seems like spotify not enabled on our side",
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
	token := base64.StdEncoding.EncodeToString([]byte(
		fmt.Sprintf(
			"%s:%s",
			neededIntegration.ClientID.String,
			neededIntegration.ClientSecret.String,
		),
	))

	resp, err := req.R().
		SetFormData(map[string]string{
			"grant_type":   "authorization_code",
			"redirect_uri": neededIntegration.RedirectURL.String,
			"code":         dto.Code,
		}).
		SetHeader("Authorization", fmt.Sprintf("Basic %s", token)).
		SetResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://accounts.spotify.com/api/token")
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "cannot get tokens")
	}
	if !resp.IsSuccess() {
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))
		return fiber.NewError(401, "seems like code is invalid")
	}

	channelIntegration.AccessToken = null.StringFrom(data.AccessToken)
	channelIntegration.RefreshToken = null.StringFrom(data.RefreshToken)

	err = services.DB.
		Save(channelIntegration).Error

	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update integration")
	}

	return nil
}

func handleGetProfile(channelId string, services types.Services) (*spotify.SpotifyProfile, error) {
	integration, err := helpers.GetIntegration(channelId, "SPOTIFY", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, err
	}

	if integration == nil {
		return nil, nil
	}

	spoty := spotify.New(integration, services.DB)
	profile, err := spoty.GetProfile()
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(400, "cannot get spotify profile")
	}

	return profile, nil
}

func handleLogout(channelId string, services types.Services) error {
	integration, err := helpers.GetIntegration(channelId, "SPOTIFY", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return err
	}
	if integration == nil {
		return fiber.NewError(http.StatusNotFound, "integration not found")
	}

	err = services.DB.Delete(&integration).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
