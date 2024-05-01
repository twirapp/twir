package commands

import (
	"context"
	"fmt"
)

func (c *CachedCommandsClient) Invalidate(ctx context.Context, channelID string) error {
	err := c.redis.Del(ctx, buildChannelCommandsCacheKey(channelID)).Err()
	if err != nil {
		return fmt.Errorf("failed to delete commands from cache: %w", err)
	}

	return nil
}
