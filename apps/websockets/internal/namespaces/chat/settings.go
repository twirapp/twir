package chat

import (
	"errors"

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

	return c.SendEvent(
		userId,
		"settings",
		map[string]any{
			"channelId":          user.ID,
			"channelName":        user.Login,
			"channelDisplayName": user.DisplayName,
			"globalBadges":       globalBadgesReq.Data.Badges,
			"channelBadges":      channelBadgesReq.Data.Badges,
		},
	)
}
