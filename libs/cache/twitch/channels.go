package twitch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kvizyx/twitchy/helix"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	[]helix.SearchChannel,
	error,
) {
	if searchString == "" {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("searchString", searchString),
	)

	if bytes, _ := c.redis.Get(
		ctx,
		buildChannelsSearchCacheKeyForId(searchString),
	).Bytes(); len(bytes) > 0 {
		var channels []helix.SearchChannel
		if err := json.Unmarshal(bytes, &channels); err != nil {
			return nil, err
		}

		return channels, nil
	}

	twitchSearchUsersReq, err := c.client.Search.SearchChannels(
		ctx, helix.SearchChannelsRequest{
			Query: searchString,
		},
	)
	if err != nil {
		return nil, err
	}

	channelsBytes, err := json.Marshal(twitchSearchUsersReq.Data)
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

	return twitchSearchUsersReq.Data, nil
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

	twitchGetChannelReq, err := c.client.Channels.GetChannelInformation(
		ctx, helix.GetChannelInformationRequest{
			BroadcasterIDs: []string{channelId},
		},
	)
	if err != nil {
		return nil, err
	}

	if len(twitchGetChannelReq.Data) == 0 {
		return nil, nil
	}

	channelBytes, err := json.Marshal(twitchGetChannelReq.Data[0])
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

	return &twitchGetChannelReq.Data[0], nil
}
