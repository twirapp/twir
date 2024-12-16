package twitch_channels

import (
	"context"

	"github.com/nicklaw5/helix/v2"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	CachedTwitchClient *twitchcahe.CachedTwitchClient
}

func New(opts Opts) *Service {
	return &Service{
		cachedTwitchClient: opts.CachedTwitchClient,
	}
}

type Service struct {
	cachedTwitchClient *twitchcahe.CachedTwitchClient
}

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
