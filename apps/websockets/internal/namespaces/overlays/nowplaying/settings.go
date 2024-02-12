package nowplaying

import (
	"github.com/satont/twir/apps/websockets/internal/protoutils"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/overlays_now_playing"
)

func (c *NowPlaying) SendSettings(userId string, overlayId string) error {
	entity := model.ChannelOverlayNowPlaying{}
	if err := c.gorm.
		Where(
			"channel_id = ? AND id = ?",
			userId,
			overlayId,
		).
		First(&entity).
		Error; err != nil {
		return err
	}

	d, err := protoutils.CreateJsonWithProto(
		&overlays_now_playing.Settings{
			Id:        entity.ID.String(),
			Preset:    entity.Preset.String(),
			ChannelId: entity.ChannelID,
		},
		nil,
	)
	if err != nil {
		return err
	}

	return c.SendEvent(
		userId,
		"settings",
		d,
	)
}
