package dudes

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

type settings struct {
	model.ChannelsOverlaysDudes
	ChannelID          string `json:"channelId"`
	ChannelName        string `json:"channelName"`
	ChannelDisplayName string `json:"channelDisplayName"`
}

func (c *Dudes) SendSettings(userId string, overlayId string) error {
	entity := model.ChannelsOverlaysDudes{}

	query := c.gorm.Where(`"channel_id" = ?`, userId)

	if overlayId != "" {
		query = query.Where("id = ?", overlayId)
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

	data := settings{
		ChannelID:             user.ID,
		ChannelName:           user.Login,
		ChannelDisplayName:    user.DisplayName,
		ChannelsOverlaysDudes: entity,
	}

	return c.SendEvent(
		userId,
		"settings",
		data,
	)
}
