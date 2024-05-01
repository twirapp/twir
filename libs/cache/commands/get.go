package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *CachedCommandsClient) GetCommands(
	ctx context.Context,
	channelID string,
) ([]model.ChannelsCommands, error) {
	if channelID == "" {
		return nil, fmt.Errorf("channelID is required")
	}

	var commands []model.ChannelsCommands
	cacheBytes, err := c.redis.Get(ctx, buildChannelCommandsCacheKey(channelID)).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get commands from cache: %w", err)
	}
	if len(cacheBytes) > 0 {
		if err := json.Unmarshal(cacheBytes, &commands); err != nil {
			return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
		}
		return commands, nil
	}

	err = c.db.
		WithContext(ctx).
		Model(&model.ChannelsCommands{}).
		Where(`channels_commands."channelId" = ? AND channels_commands."enabled" = ?`, channelID, true).
		Joins("Group").
		Preload("Responses").
		WithContext(ctx).
		Find(&commands).Error
	if err != nil {
		return nil, err
	}

	cacheBytes, err = json.Marshal(commands)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal commands: %w", err)
	}

	if err := c.redis.Set(
		ctx,
		buildChannelCommandsCacheKey(channelID),
		cacheBytes,
		channelCommandsCacheTTL,
	).Err(); err != nil {
		return nil, fmt.Errorf("failed to set commands to cache: %w", err)
	}

	return commands, nil
}
