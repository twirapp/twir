package handlers

import (
	"context"
	"fmt"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
	"go.uber.org/zap"
)

func (c *Handlers) OnUserJoin(message irc.UserJoinMessage) {
	ctx := context.Background()

	stream := &model.ChannelsStreams{}
	err := c.db.
		WithContext(ctx).
		Where(`"userLogin" = ?`, message.Channel).
		Find(stream).Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	if stream.ID == "" {
		return
	}

	ignoredUser := &model.IgnoredUser{}
	err = c.db.
		WithContext(ctx).
		Where("login = ?", message.User).
		Find(ignoredUser).Error
	if err != nil {
		zap.S().Error(err)
		return
	}

	if ignoredUser.ID != "" {
		return
	}

	redisKey := fmt.Sprintf("events:first-stream-user-join:%s:%s", stream.ID, message.User)

	res, err := c.redis.Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		zap.S().Error(err)
		return
	}

	if res != "" {
		return
	}

	_, err = c.eventsGrpc.StreamFirstUserJoin(
		ctx, &events.StreamFirstUserJoinMessage{
			BaseInfo: &events.BaseInfo{
				ChannelId: stream.UserId,
			},
			UserName: message.User,
		},
	)

	if err != nil {
		zap.S().Error(err)
	}

	_, err = c.redis.Set(
		ctx,
		redisKey,
		message.User,
		49*time.Hour,
	).Result()

	if err != nil {
		zap.S().Error(err)
	}
}
