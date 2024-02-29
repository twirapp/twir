package public_settings

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/satont/twir/apps/api/internal/converters"
	"github.com/satont/twir/apps/api/internal/helpers"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/channels_public_settings"
	"gorm.io/gorm"
)

type PublicSettings struct {
	*impl_deps.Deps
}

func (c *PublicSettings) ChannelsPublicSettingsUpdate(
	ctx context.Context,
	req *channels_public_settings.UpdateRequest,
) (*channels_public_settings.Settings, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	currentSettings := &model.ChannelPublicSettings{}
	if err := c.Deps.Db.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			dashboardId,
		).
		Preload("SocialLinks").
		// init default settings
		FirstOrInit(
			currentSettings,
			&model.ChannelPublicSettings{
				ChannelID: dashboardId,
			},
		).
		Error; err != nil {
		return nil, err
	}

	links := make([]model.ChannelPublicSettingsSocialLink, 0, len(req.SocialLinks))
	for _, link := range req.SocialLinks {
		links = append(
			links,
			model.ChannelPublicSettingsSocialLink{
				ID:         uuid.New(),
				SettingsID: currentSettings.ID,
				Title:      link.Title,
				Href:       link.Href,
			},
		)
	}

	err = c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			for _, l := range currentSettings.SocialLinks {
				if err := tx.Delete(l).Error; err != nil {
					return err
				}
			}

			currentSettings.Description = null.StringFromPtr(req.Description)
			currentSettings.SocialLinks = links

			return c.Db.WithContext(ctx).Save(&currentSettings).Error
		},
	)

	if err != nil {
		return nil, err
	}

	return converters.ChannelsPublicSettingsModelToRpc(currentSettings), nil
}
