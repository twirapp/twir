package kappagen

import (
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

func (c *Kappagen) SendSettings(userId string) error {
	settings := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, userId, "kappagen_overlay").
		Find(settings).
		Error
	if err != nil {
		return err
	}

	if settings.ID == "" {
		return nil
	}

	twitchClient, err := twitch.NewUserClient(userId, c.config, c.tokensGrpc)
	if err != nil {
		return err
	}

	usersReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{userId},
		},
	)
	if err != nil {
		return err
	}

	if len(usersReq.Data.Users) == 0 {
		return errors.New("cannot get user")
	}

	user := usersReq.Data.Users[0]

	settingsMap := make(map[string]any)
	settingsMap["channelId"] = user.ID
	settingsMap["channelName"] = user.Login

	dbSettingsMap := make(map[string]any)
	if err := json.Unmarshal(settings.Settings, &dbSettingsMap); err != nil {
		return fmt.Errorf("cannot unmarshal dbSettings: %w", err)
	}

	for key, value := range dbSettingsMap {
		settingsMap[key] = value
	}

	return c.SendEvent(
		userId,
		"settings",
		settingsMap,
	)
}
