package spotify

import (
	"encoding/base64"
	"fmt"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/helpers"
	"io"
	"net/http"
	"net/url"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/satont/tsuwari/libs/integrations/spotify"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (c *Spotify) getAuthLinkService() (*string, error) {
	integration := model.Integrations{}
	err := c.services.Gorm.Where(`"service" = ?`, "SPOTIFY").First(&integration).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(
			404,
			"spotify not enabled on our side. Please be patient.",
		)
	}

	parsedUrl, _ := url.Parse("https://accounts.spotify.com/authorize")

	q := parsedUrl.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", "user-read-currently-playing")
	q.Add("redirect_uri", integration.RedirectURL.String)
	parsedUrl.RawQuery = q.Encode()

	str := parsedUrl.String()

	return &str, nil
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

func (c *Spotify) postService(channelId string, dto *tokenDto) error {
	channelIntegration, err := helpers.GetIntegration(channelId, "SPOTIFY", c.services.Gorm)
	if err != nil {
		c.services.Logger.Error(err)
		return err
	}

	neededIntegration := model.Integrations{}
	err = c.services.Gorm.
		Where("service = ?", "SPOTIFY").
		First(&neededIntegration).
		Error
	if err != nil {
		c.services.Logger.Error(err)
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
			Integration:   &neededIntegration,
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
		data, _ := io.ReadAll(resp.Body)
		fmt.Println(string(data))
		return fiber.NewError(401, "seems like code is invalid")
	}

	channelIntegration.AccessToken = null.StringFrom(data.AccessToken)
	channelIntegration.RefreshToken = null.StringFrom(data.RefreshToken)

	err = c.services.Gorm.
		Save(channelIntegration).Error

	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update integration")
	}

	return nil
}

func (c *Spotify) getService(channelId string) (*spotify.SpotifyProfile, error) {
	integration, err := helpers.GetIntegration(channelId, "SPOTIFY", c.services.Gorm)
	if err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	if integration == nil {
		return nil, nil
	}

	spoty := spotify.New(integration, c.services.Gorm)
	profile, err := spoty.GetProfile()
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(400, "cannot get spotify profile")
	}

	return profile, nil
}

func (c *Spotify) logoutService(channelId string) error {
	integration, err := helpers.GetIntegration(channelId, "SPOTIFY", c.services.Gorm)
	if err != nil {
		c.services.Logger.Error(err)
		return err
	}
	if integration == nil {
		return fiber.NewError(http.StatusNotFound, "integration not found")
	}

	err = c.services.Gorm.Delete(&integration).Error
	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
