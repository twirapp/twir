package public_settings

import (
	"context"

	"github.com/satont/twir/apps/api/internal/converters"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/channels_public_settings"
)

type PublicSettings struct {
	*impl_deps.Deps
}

func (c *PublicSettings) GetPublicSettings(
	ctx context.Context,
	req *channels_public_settings.GetRequest,
) (*channels_public_settings.Settings, error) {
	entity := &model.ChannelPublicSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			req.ChannelId,
		).
		Preload("SocialLinks").
		Find(entity).Error; err != nil {
		return nil, err
	}

	return converters.ChannelsPublicSettingsModelToRpc(entity), nil
}
