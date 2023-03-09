package lastfm

import (
	"fmt"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/channels/integrations/helpers"
	"net/http"

	"gorm.io/gorm"

	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	lfm "github.com/shkh/lastfm-go/lastfm"
)

type Lastfm struct {
	model.ChannelsIntegrations
	Data map[string]any `json:"data"`
}

func (c *LastFM) postService(channelId string, dto *lastfmDto) error {
	integration, err := helpers.GetIntegration(channelId, "LASTFM", c.services.Gorm)
	if err != nil {
		c.services.Logger.Error(err)
		return err
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = c.services.Gorm.
			Where("service = ?", "LASTFM").
			First(&neededIntegration).
			Error
		if err != nil {
			c.services.Logger.Error(err)
			return fiber.NewError(
				http.StatusInternalServerError,
				"seems like lastfm not enabled on our side",
			)
		}

		integration = &model.ChannelsIntegrations{
			ID:            uuid.NewV4().String(),
			ChannelID:     channelId,
			IntegrationID: neededIntegration.ID,
			Integration:   &neededIntegration,
		}
	}

	integration.Enabled = true

	api := lfm.New(
		integration.Integration.APIKey.String,
		integration.Integration.ClientSecret.String,
	)
	err = api.LoginWithToken(dto.Code)
	sessionKey := api.GetSessionKey()

	integration.APIKey = null.StringFrom(sessionKey)

	if err = c.services.Gorm.Save(integration).Error; err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update faceit data")
	}

	return nil
}

func (c *LastFM) authService() (string, error) {
	// http://www.last.fm/api/auth/?api_key=xxx&cb=http://example.com
	neededIntegration := model.Integrations{}
	err := c.services.Gorm.
		Where("service = ?", "LASTFM").
		First(&neededIntegration).
		Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return "", fiber.NewError(
			404,
			"lastfm not enabled on our side. Please be patient.",
		)
	}

	return fmt.Sprintf(
		"http://www.last.fm/api/auth/?api_key=%s&cb=%s",
		neededIntegration.APIKey.String,
		neededIntegration.RedirectURL.String,
	), nil
}

func (c *LastFM) getService(channelId string) (*LastfmProfile, error) {
	integration, err := helpers.GetIntegration(channelId, "LASTFM", c.services.Gorm)
	if err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	if integration == nil {
		return nil, nil
	}

	api := lfm.New(
		integration.Integration.APIKey.String,
		integration.Integration.ClientSecret.String,
	)
	api.SetSession(integration.APIKey.String)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}
	info, err := api.User.GetInfo(make(map[string]interface{}))
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return &LastfmProfile{
		Name:      info.Name,
		Image:     info.Images[len(info.Images)-1].Url,
		PlayCount: info.PlayCount,
	}, nil
}

func (c *LastFM) logoutService(channelId string) error {
	integration, err := helpers.GetIntegration(channelId, "LASTFM", c.services.Gorm)
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
