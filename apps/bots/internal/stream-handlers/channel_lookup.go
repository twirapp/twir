package stream_handlers

import (
	"context"
	"errors"

	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
)

type twitchChannelLookupResult struct {
	ID string
}

func (c *PubSubHandlers) findTwitchChannelByPlatformChannelID(
	ctx context.Context,
	platformChannelID string,
) (twitchChannelLookupResult, bool, error) {
	channel, err := c.channelService.GetChannelByPlatformChannelID(
		ctx,
		platform.PlatformTwitch,
		platformChannelID,
	)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return twitchChannelLookupResult{}, false, nil
		}

		return twitchChannelLookupResult{}, false, err
	}

	return twitchChannelLookupResult{ID: channel.ID.String()}, true, nil
}
