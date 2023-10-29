package chat

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/olahol/melody"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

func (c *Chat) sendSettings(session *melody.Session) error {
	userId, exists := session.Get("userId")
	if !exists {
		return nil
	}

	channelId := userId.(string)

	settings := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, channelId, "chat_overlay").
		Find(settings).
		Error
	if err != nil {
		return err
	}

	twitchClient, err := twitch.NewUserClient(channelId, c.config, c.tokensGrpc)
	if err != nil {
		return err
	}

	usersReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{channelId},
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
			BroadcasterID: channelId,
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
		userId.(string),
		"settings",
		map[string]any{
			"channelId":          channelId,
			"channelName":        user.Login,
			"channelDisplayName": user.DisplayName,
			"globalBadges":       globalBadgesReq.Data.Badges,
			"channelBadges":      channelBadgesReq.Data.Badges,
		},
	)
}
