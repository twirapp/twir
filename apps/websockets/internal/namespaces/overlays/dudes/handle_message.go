package dudes

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/olahol/melody"
	"github.com/samber/lo"
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
	emptySettings := overlays.DudesUserSettings{
		UserID: userId,
	}

	twitchClient, err := twitch.NewAppClient(c.config, c.tokensGrpc)
	if err != nil {
		c.logger.Error(err.Error())
		return c.SendEvent(
			channelId,
			"userSettings",
			&emptySettings,
		)
	}

	usersReq, err := retry.DoWithData(
		func() (*helix.UsersResponse, error) {
			return twitchClient.GetUsers(
				&helix.UsersParams{
					IDs: []string{userId},
				},
			)
		},
		retry.Attempts(5),
	)
	if err != nil {
		c.logger.Error(err.Error())
		return c.SendEvent(
			channelId,
			"userSettings",
			&emptySettings,
		)
	}
	if usersReq.ErrorMessage != "" {
		return c.SendEvent(
			channelId,
			"userSettings",
			&emptySettings,
		)
	}
	if len(usersReq.Data.Users) == 0 {
		c.logger.Warn("cannot get user")
		return c.SendEvent(
			channelId,
			"userSettings",
			&emptySettings,
		)
	}

	user := usersReq.Data.Users[0]
	emptySettings.UserDisplayName = &user.DisplayName
	emptySettings.UserName = &user.Login

	err = c.gorm.
		Where("channel_id = ? AND user_id = ?", channelId, userId).
		First(&entity).Error
	if err != nil {
		c.logger.Error(err.Error())
		return c.SendEvent(
			channelId,
			"userSettings",
			&emptySettings,
		)
	}

	var sprite *overlays.DudesSprite
	if entity.DudeSprite != nil {
		sprite = lo.ToPtr(overlays.DudesSprite(*entity.DudeSprite))
	}

	c.SendEvent(
		channelId,
		"userSettings",
		&overlays.DudesUserSettings{
			DudeColor:       entity.DudeColor,
			DudeSprite:      sprite,
			UserID:          user.ID,
			UserDisplayName: &user.DisplayName,
			UserName:        &user.Login,
		},
	)

	return nil
}
