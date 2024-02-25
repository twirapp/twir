package dudes

import (
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/types/types/overlays"
)

func (c *Dudes) handleMessage(session *melody.Session, msg []byte) {
	channelId, ok := session.Get("userId")
	if channelId == nil || channelId == "" || !ok {
		return
	}

	var overlayId string
	id, ok := session.Get("id")

	if id != nil || ok {
		casted, castOk := id.(string)
		if castOk {
			overlayId = casted
		}
	}

	data := &types.WebSocketMessage{
		CreatedAt: time.Now().UTC().String(),
	}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	if data.EventName == "getSettings" {
		err := c.SendSettings(channelId.(string), overlayId)
		if err != nil {
			c.logger.Error(err.Error())
		}
	}

	if data.EventName == "getUserSettings" {
		userId, ok := data.Data.(string)
		if !ok {
			return
		}

		err := c.SendUserSettings(channelId.(string), userId)
		if err != nil {
			c.logger.Error(err.Error())
		}
	}
}

func (c *Dudes) SendUserSettings(
	channelId string,
	userId string,
) error {
	entity := model.ChannelsOverlaysDudesUserSettings{}

	err := c.gorm.
		Where("channel_id = ? AND user_id = ?", channelId, userId).
		First(&entity).Error
	if err != nil {
		return err
	}

	twitchClient, err := twitch.NewAppClient(c.config, c.tokensGrpc)
	if err != nil {
		c.logger.Error(err.Error())
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
	if usersReq.ErrorMessage != "" {
		return errors.New(usersReq.ErrorMessage)
	}
	if len(usersReq.Data.Users) == 0 {
		return nil
	}

	user := usersReq.Data.Users[0]

	c.SendEvent(
		channelId,
		"userSettings",
		&overlays.DudesUserSettings{
			DudeColor:  entity.DudeColor,
			DudeSprite: entity.DudeSprite,
			UserID:     user.ID,
			UserName:   user.DisplayName,
			UserLogin:  user.Login,
		},
	)

	return nil
}
