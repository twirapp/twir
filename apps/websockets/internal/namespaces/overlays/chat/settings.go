package chat

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"gorm.io/gorm"
)

type settings struct {
	model.ChatOverlaySettings
	ChannelID          string            `json:"channel_id"`
	ChannelName        string            `json:"channel_name"`
	ChannelDisplayName string            `json:"channel_display_name"`
	GlobalBadges       []helix.ChatBadge `json:"global_badges"`
	ChannelBadges      []helix.ChatBadge `json:"channel_badges"`
}

func (c *Chat) SendSettings(userId string, overlayId string) error {
	entity := model.ChatOverlaySettings{}

	query := c.gorm.Where(`"channel_id" = ?`, userId)

	if overlayId != "" {
		query = query.Where("id = ?", overlayId)
	}

	err := query.Order("created_at asc").First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

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

	data := settings{
		ChannelID:           user.ID,
		ChannelName:         user.Login,
		ChannelDisplayName:  user.DisplayName,
		GlobalBadges:        globalBadgesReq.Data.Badges,
		ChannelBadges:       channelBadgesReq.Data.Badges,
		ChatOverlaySettings: entity,
	}

	return c.SendEvent(
		userId,
		"settings",
		data,
	)
}
