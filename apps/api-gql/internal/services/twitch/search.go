package twitch

import (
	"context"

	"github.com/nicklaw5/helix/v2"
)

func (c *Service) SearchByName(ctx context.Context, query string) ([]helix.Channel, error) {
	if query == "" {
		return nil, nil
	}

	channels, err := c.cachedTwitchClient.SearchChannels(ctx, query)
	if err != nil {
		return nil, err
	}

	return channels, err
}
