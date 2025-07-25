package dudes

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/websockets/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/types/types/overlays"
)

func (c *Dudes) handleMessage(session *melody.Session, msg []byte) {
	channelId, ok := session.Get("userId")
	if channelId == nil || channelId == "" || !ok {
		return
	}

	data := &types.WebSocketMessage{
		CreatedAt: time.Now().UTC().String(),
	}
	err := json.Unmarshal(msg, data)
	if err != nil {
		c.logger.Error(err.Error())
		return
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

	err := c.gorm.
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
			DudeColor:  entity.DudeColor,
			DudeSprite: sprite,
			UserID:     userId,
		},
	)

	return nil
}
