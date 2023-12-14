package kappagen

import (
	"errors"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

type settings struct {
	model.KappagenOverlaySettings
	ChannelID   string `json:"channelId"`
	ChannelName string `json:"channelName"`
}

func (c *Kappagen) SendSettings(userId string) error {
	entity := &model.ChannelModulesSettings{}
	err := c.gorm.
		Where(`"channelId" = ? AND "type" = ?`, userId, "kappagen_overlay").
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

	parsedEntitySettings := model.KappagenOverlaySettings{}
	err = json.Unmarshal(entity.Settings, &parsedEntitySettings)
	if err != nil {
		return err
	}

	settingsMap := settings{
		ChannelID:               user.ID,
		ChannelName:             user.Login,
		KappagenOverlaySettings: parsedEntitySettings,
	}

	return c.SendEvent(
		userId,
		"settings",
		settingsMap,
	)
}
