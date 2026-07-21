package stream_handlers

import (
	"context"
	"errors"

	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
)

type twitchChannelLookupResult struct {
	ID    string
	BotID string
}

func (c *PubSubHandlers) findTwitchChannelByPlatformUserID(
	ctx context.Context,
	platformUserID string,
) (twitchChannelLookupResult, bool, error) {
	channel, err := c.channelService.GetChannelByPlatformUserID(
		ctx,
		platformUserID,
		platform.PlatformTwitch,
	)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return twitchChannelLookupResult{}, false, nil
		}

		return twitchChannelLookupResult{}, false, err
	}

	return twitchChannelLookupResult{ID: channel.ID.String(), BotID: channel.BotID}, true, nil
}
