package be_right_back

import (
	"fmt"

	"github.com/goccy/go-json"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *BeRightBack) SendSettings(userId string) error {
	settings := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, userId, "be_right_back_overlay").
		Find(settings).
		Error
	if err != nil {
		return err
	}

	if settings.ID == "" {
		return nil
	}

	data := model.ChannelModulesSettingsBeRightBack{}
	if err := json.Unmarshal(settings.Settings, &data); err != nil {
		return fmt.Errorf("cannot unmarshal dbSettings: %w", err)
	}

	return c.SendEvent(
		userId,
		"settings",
		data,
	)
}
