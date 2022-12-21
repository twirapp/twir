package lastfm

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
	lfm "github.com/shkh/lastfm-go/lastfm"
)

type Lastfm struct {
	model.ChannelsIntegrations
	Data map[string]any `json:"data"`
}

func handlePost(channelId string, dto *lastfmDto, services types.Services) error {
	integration, err := helpers.GetIntegration(channelId, "LASTFM", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return err
	}

	if integration == nil {
		neededIntegration := model.Integrations{}
		err = services.DB.
			Where("service = ?", "LASTFM").
			First(&neededIntegration).
			Error
		if err != nil {
			services.Logger.Sugar().Error(err)
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
	err = api.LoginWithToken(dto.Token)
	sessionKey := api.GetSessionKey()

	integration.APIKey = null.StringFrom(sessionKey)

	if err = services.DB.Save(integration).Error; err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot update faceit data")
	}

	return nil
}

func handleAuth(services types.Services) (string, error) {
	// http://www.last.fm/api/auth/?api_key=xxx&cb=http://example.com
	neededIntegration := model.Integrations{}
	err := services.DB.
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

func handleProfile(channelId string, services types.Services) (*LastfmProfile, error) {
	integration, err := helpers.GetIntegration(channelId, "LASTFM", services.DB)
	if err != nil {
		services.Logger.Sugar().Error(err)
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
		spew.Dump(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}
	info, err := api.User.GetInfo(make(map[string]interface{}))
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return &LastfmProfile{
		Name:      info.Name,
		Image:     info.Images[len(info.Images)-1].Url,
		PlayCount: info.PlayCount,
	}, nil
}

func handleLogout(channelId string, services types.Services) error {
	integration, err := helpers.GetIntegration(channelId, "LASTFM", services.DB)
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
