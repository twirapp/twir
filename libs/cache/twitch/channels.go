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

const channelsByIdCacheKey = "cache:twir:twitch:channels:byId:"
const channelsByIdTTL = 10 * time.Second

func buildChannelsByIdCacheKeyForId(channelId string) string {
	return channelsByIdCacheKey + channelId
}

func (c *CachedTwitchClient) GetChannelInformationById(
	ctx context.Context,
	channelId string,
) (
	*helix.ChannelInformation,
	error,
) {
	if channelId == "" {
		return nil, nil
	}

	if bytes, _ := c.redis.Get(
		ctx,
		buildChannelsByIdCacheKeyForId(channelId),
	).Bytes(); len(bytes) > 0 {
		var channel helix.ChannelInformation
		if err := json.Unmarshal(bytes, &channel); err != nil {
			return nil, err
		}

		return &channel, nil
	}

	twitchGetChannelReq, err := c.client.GetChannelInformation(
		&helix.GetChannelInformationParams{
			BroadcasterIDs: []string{channelId},
		},
	)
	if err != nil {
		return nil, err
	}
	if twitchGetChannelReq.ErrorMessage != "" {
		return nil, err
	}

	if len(twitchGetChannelReq.Data.Channels) == 0 {
		return nil, nil
	}

	channelBytes, err := json.Marshal(twitchGetChannelReq.Data.Channels[0])
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(
		ctx,
		buildChannelsByIdCacheKeyForId(channelId),
		channelBytes,
		channelsByIdTTL,
	).Err(); err != nil {
		return nil, err
	}

	return &twitchGetChannelReq.Data.Channels[0], nil
}
