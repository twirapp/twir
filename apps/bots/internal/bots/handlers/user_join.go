package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
)

type OnUserJoinOpts struct {
	Channel string
	User    string
}

func (c *Handlers) OnUserJoin(message OnUserJoinOpts) {
	ctx := context.Background()

	stream := &model.ChannelsStreams{}
	err := c.db.
		WithContext(ctx).
		Where(`"userLogin" = ?`, message.Channel).
		Find(stream).Error
	if err != nil {
		c.logger.Error(
			"cannot get channel stream",
			slog.Any("err", err),
			slog.String(
				"login",
				message.Channel,
			),
		)
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
		c.logger.Error(
			"cannot get ignored user",
			slog.Any("err", err),
			slog.String(
				"login",
				message.User,
			),
		)
		return
	}

	if ignoredUser.ID != "" {
		return
	}

	redisKey := fmt.Sprintf("events:first-stream-user-join:%s:%s", stream.ID, message.User)

	res, err := c.redis.Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		c.logger.Error(
			"cannot first join",
			slog.Any("err", err),
			slog.String("login", message.User),
		)
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
		c.logger.Error(
			"cannot fire first join to events",
			slog.Any("err", err),
			slog.String("login", message.User),
			slog.String("channelId", stream.UserId),
		)
	}

	_, err = c.redis.Set(
		ctx,
		redisKey,
		message.User,
		49*time.Hour,
	).Result()

	if err != nil {
		c.logger.Error(
			"cannot set first join to redis",
			slog.Any("err", err),
			slog.String("login", message.User),
			slog.String("channelId", stream.UserId),
		)
	}
}
