package chat_client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/events"
)

func (c *ChatClient) onUserJoin(message irc.UserJoinMessage) {
	ctx := context.Background()

	stream := &model.ChannelsStreams{}
	err := c.services.DB.
		WithContext(ctx).
		Where(`"userLogin" = ?`, message.Channel).
		Find(stream).Error
	if err != nil {
		c.services.Logger.Error(
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
	err = c.services.DB.
		WithContext(ctx).
		Where("login = ?", message.User).
		Find(ignoredUser).Error
	if err != nil {
		c.services.Logger.Error(
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

	res, err := c.services.Redis.Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		c.services.Logger.Error(
			"cannot first join",
			slog.Any("err", err),
			slog.String("login", message.User),
		)
		return
	}

	if res != "" {
		return
	}

	_, err = c.services.EventsGrpc.StreamFirstUserJoin(
		ctx, &events.StreamFirstUserJoinMessage{
			BaseInfo: &events.BaseInfo{
				ChannelId: stream.UserId,
			},
			UserName: message.User,
		},
	)
	if err != nil {
		c.services.Logger.Error(
			"cannot fire first join to events",
			slog.Any("err", err),
			slog.String("login", message.User),
			slog.String("channelId", stream.UserId),
		)
	}

	_, err = c.services.Redis.Set(
		ctx,
		redisKey,
		message.User,
		49*time.Hour,
	).Result()

	if err != nil {
		c.services.Logger.Error(
			"cannot set first join to redis",
			slog.Any("err", err),
			slog.String("login", message.User),
			slog.String("channelId", stream.UserId),
		)
	}
}
