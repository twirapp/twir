package messagehandler

import (
	"context"
	"time"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
)

func (c *MessageHandler) handleFirstStreamUserJoin(ctx context.Context, msg handleMessage) error {
	if msg.DbStream == nil {
		return nil
	}

	ignoredUser := &model.IgnoredUser{}
	err := c.gorm.
		WithContext(ctx).
		Where("login = ? OR id = ?", msg.ChatterUserLogin, msg.ChatterUserId).
		Find(ignoredUser).Error
	if err != nil {
		return err
	}
	if ignoredUser.ID != "" {
		return nil
	}

	redisKey := "first:stream:user:join:" + msg.ChatterUserId

	exists, err := c.redis.Exists(ctx, redisKey).Result()
	if err != nil {
		return err
	}

	if exists == 1 {
		return nil
	}

	err = c.redis.Set(ctx, redisKey, "", 48*time.Hour).Err()
	if err != nil {
		return err
	}

	_, err = c.eventsGrpc.StreamFirstUserJoin(
		ctx, &events.StreamFirstUserJoinMessage{
			BaseInfo: &events.BaseInfo{
				ChannelId: msg.BroadcasterUserId,
			},
			UserName: msg.ChatterUserName,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
