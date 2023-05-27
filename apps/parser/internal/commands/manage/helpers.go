package manage

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func dropRedisCache(
	ctx context.Context,
	redis *redis.Client,
	logger *zap.SugaredLogger,
	channelId string,
) {
	err := redis.Del(
		ctx,
		fmt.Sprintf("fiber:cache:/v1/channels/%s/commands_GET", channelId),
	).Err()

	if err != nil {
		logger.Error(err)
	}

	err = redis.Del(
		ctx,
		fmt.Sprintf("fiber:cache:/v1/channels/%s/commands_GET_body", channelId),
	).Err()

	if err != nil {
		logger.Error(err)
	}
}
