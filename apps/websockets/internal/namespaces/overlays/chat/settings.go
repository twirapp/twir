package chat

import (
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

func (c *Chat) SendSettings(userId string) error {
	settings := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, userId, "chat_overlay").
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

	channelBadgesReq, err := twitchClient.GetChannelChatBadges(
		&helix.GetChatBadgeParams{
			BroadcasterID: userId,
		},
	)
	if err != nil {
		return err
	}

	globalBadgesReq, err := twitchClient.GetGlobalChatBadges()
	if err != nil {
		return err
	}

	settingsMap := make(map[string]any)
	settingsMap["channelId"] = user.ID
	settingsMap["channelName"] = user.Login
	settingsMap["channelDisplayName"] = user.DisplayName
	settingsMap["globalBadges"] = globalBadgesReq.Data.Badges
	settingsMap["channelBadges"] = channelBadgesReq.Data.Badges

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
