package twitch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nicklaw5/helix/v2"
)

const channelsSearchKey = "cache:twir:twitch:channels:search:"
const channelsSearchTTL = 5 * time.Hour

func buildChannelsSearchCacheKeyForId(searchString string) string {
	return channelsSearchKey + searchString
}

func (c *CachedTwitchClient) SearchChannels(
	ctx context.Context,
	searchString string,
) (
	[]helix.Channel,
	error,
) {
	if searchString == "" {
		return nil, nil
	}

	if bytes, _ := c.redis.Get(
		ctx,
		buildChannelsSearchCacheKeyForId(searchString),
	).Bytes(); len(bytes) > 0 {
		var channels []helix.Channel
		if err := json.Unmarshal(bytes, &channels); err != nil {
			return nil, err
		}

		return channels, nil
	}

	twitchSearchUsersReq, err := c.client.SearchChannels(
		&helix.SearchChannelsParams{
			Channel: searchString,
		},
	)
	if err != nil {
		return nil, err
	}
	if twitchSearchUsersReq.ErrorMessage != "" {
		return nil, err
	}

	channelsBytes, err := json.Marshal(twitchSearchUsersReq.Data.Channels)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(
		ctx,
		buildChannelsSearchCacheKeyForId(searchString),
		channelsBytes,
		channelsSearchTTL,
	).Err(); err != nil {
		return nil, err
	}

	return twitchSearchUsersReq.Data.Channels, nil
}
