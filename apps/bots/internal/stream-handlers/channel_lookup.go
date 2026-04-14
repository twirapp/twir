package stream_handlers

import (
	"context"
	"errors"

	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	user_platform_accounts "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
)

type twitchChannelLookupResult struct {
	ID    string
	BotID string
}

func (c *PubSubHandlers) findTwitchChannelByPlatformUserID(
	ctx context.Context,
	platformUserID string,
) (twitchChannelLookupResult, bool, error) {
	account, err := c.userPlatformAccountsRepo.GetByPlatformUserID(ctx, platform.PlatformTwitch, platformUserID)
	if err != nil {
		if errors.Is(err, user_platform_accounts.ErrNotFound) {
			return twitchChannelLookupResult{}, false, nil
		}

		return twitchChannelLookupResult{}, false, err
	}

	channel, err := c.channelsRepo.GetByUserIDAndPlatform(ctx, account.UserID, platform.PlatformTwitch)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return twitchChannelLookupResult{}, false, nil
		}

		return twitchChannelLookupResult{}, false, err
	}

	return twitchChannelLookupResult{ID: channel.ID, BotID: channel.BotID}, true, nil
}
