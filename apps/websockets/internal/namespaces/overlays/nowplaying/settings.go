package nowplaying

import (
	model "github.com/satont/twir/libs/gomodels"
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

	return c.SendEvent(
		userId,
		"settings",
		entity,
	)
}
