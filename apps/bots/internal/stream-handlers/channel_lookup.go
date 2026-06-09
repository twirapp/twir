package stream_handlers

import (
	"context"
	"errors"

	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type twitchChannelLookupResult struct {
	ID    string
	BotID string
}

func (c *PubSubHandlers) findTwitchChannelByPlatformUserID(
	ctx context.Context,
	platformUserID string,
) (twitchChannelLookupResult, bool, error) {
	user, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, platformUserID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return twitchChannelLookupResult{}, false, nil
		}

		return twitchChannelLookupResult{}, false, err
	}

	channel, err := c.channelsRepo.GetByTwitchUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return twitchChannelLookupResult{}, false, nil
		}

		return twitchChannelLookupResult{}, false, err
	}

	return twitchChannelLookupResult{ID: channel.ID.String(), BotID: channel.BotID}, true, nil
}
