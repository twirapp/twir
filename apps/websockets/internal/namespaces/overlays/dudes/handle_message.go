package dudes

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/websockets/types"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	twirlogger "github.com/twirapp/twir/libs/logger"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"github.com/twirapp/twir/libs/types/types/overlays"
	"gorm.io/gorm"
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

		err := c.SendUserSettings(context.Background(), channelId.(string), userId)
		if err != nil {
			c.logger.Error(err.Error())
		}
	}
}

func (c *Dudes) SendUserSettings(
	ctx context.Context,
	channelId string,
	userId string,
) error {
	emptySettings := overlays.DudesUserSettings{UserID: userId}

	queryUserID := userId
	if _, err := uuid.Parse(userId); err != nil {
		user, err := c.usersRepository.GetByPlatformID(ctx, platform.PlatformTwitch, userId)
		if err != nil {
			if !errors.Is(err, usersmodel.ErrNotFound) {
				c.logger.Error("cannot resolve user by twitch id", twirlogger.Error(err))
			}
			return c.SendEvent(channelId, "userSettings", &emptySettings)
		}
		queryUserID = user.ID.String()
	}

	entity := model.ChannelsOverlaysDudesUserSettings{}

	err := c.gorm.
		Where("channel_id = ? AND user_id = ?", channelId, queryUserID).
		First(&entity).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.logger.Error(err.Error())
		}
		return c.SendEvent(channelId, "userSettings", &emptySettings)
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
