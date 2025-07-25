package be_right_back

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/websockets/internal/protoutils"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/overlays_be_right_back"
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

	overlaySettings := overlays_be_right_back.Settings{
		Text: data.Text,
		Late: &overlays_be_right_back.Settings_Late{
			Enabled:        data.Late.Enabled,
			Text:           data.Late.Text,
			DisplayBrbTime: data.Late.DisplayBrbTime,
		},
		BackgroundColor: data.BackgroundColor,
		FontSize:        data.FontSize,
		FontColor:       data.FontColor,
		FontFamily:      data.FontFamily,
	}

	d, err := protoutils.CreateJsonWithProto(&overlaySettings, nil)
	if err != nil {
		return err
	}

	return c.SendEvent(
		userId,
		"settings",
		d,
	)
}
