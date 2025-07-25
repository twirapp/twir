package messagehandler

import (
	"context"
	"fmt"
	"time"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func (c *MessageHandler) handleFirstStreamUserJoin(ctx context.Context, msg handleMessage) error {
	if msg.EnrichedData.ChannelStream == nil {
		return nil
	}

	redisKey := fmt.Sprintf(
		"first:stream:user:join:%s:%s",
		msg.EnrichedData.ChannelStream.ID,
		msg.ChatterUserId,
	)
	exists, err := c.redis.Exists(ctx, redisKey).Result()
	if err != nil {
		return err
	}

	if exists == 1 {
		return nil
	}

	ignoredUser := &model.IgnoredUser{}
	err = c.gorm.
		WithContext(ctx).
		Where("login = ? OR id = ?", msg.ChatterUserLogin, msg.ChatterUserId).
		Find(ignoredUser).Error
	if err != nil {
		return err
	}
	if ignoredUser.ID != "" {
		return nil
	}

	err = c.redis.Set(ctx, redisKey, "", 48*time.Hour).Err()
	if err != nil {
		return err
	}

	err = c.twirBus.Events.StreamFirstUserJoin.Publish(
		ctx,
		events.StreamFirstUserJoinMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   msg.BroadcasterUserId,
				ChannelName: msg.BroadcasterUserLogin,
			},
			UserLogin: msg.ChatterUserName,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
