package handle_message

import (
	"context"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func (c *Service) getCachedChannel(ctx context.Context, channelID string) (
	channelsmodel.Channel,
	error,
) {
	cacheKey := "cache:chat_translations:channels:" + channelID
	cachedBytes, err := c.redis.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return channelsmodel.Channel{}, err
	}

	if len(cachedBytes) > 0 {
		channel := channelsmodel.Channel{}
		if err := json.Unmarshal(cachedBytes, &channel); err != nil {
			return channelsmodel.Channel{}, err
		}

		return channel, nil
	}

	channel, err := c.channelsRepository.GetByID(ctx, channelID)
	if err != nil {
		return channelsmodel.Channel{}, err
	}

	channelBytes, err := json.Marshal(channel)
	if err != nil {
		return channelsmodel.Channel{}, err
	}

	if err := c.redis.Set(
		ctx,
		cacheKey,
		channelBytes,
		5*time.Minute,
	).Err(); err != nil {
		return channelsmodel.Channel{}, err
	}

	return channel, nil
}
