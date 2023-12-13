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
	entity := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, userId, "chat_overlay").
		Find(entity).
		Error
	if err != nil {
		return err
	}

	if entity.ID == "" {
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

	var entitySettings model.ChatOverlaySettings
	if err := json.Unmarshal(entity.Settings, &entitySettings); err != nil {
		return fmt.Errorf("cannot unmarshal entitySettings: %w", err)
	}

	data := settings{
		ChannelID:           user.ID,
		ChannelName:         user.Login,
		ChannelDisplayName:  user.DisplayName,
		GlobalBadges:        globalBadgesReq.Data.Badges,
		ChannelBadges:       channelBadgesReq.Data.Badges,
		ChatOverlaySettings: entitySettings,
	}

	return c.SendEvent(
		userId,
		"settings",
		data,
	)
}
