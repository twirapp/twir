package channelservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
)

func NewChannelService(
	repo channelsrepo.Repository,
	twirbus *buscore.Bus,
	cfg config.Config,
	cache kv.KV,
	streamsrepo streamsrepository.Repository,
) *ChannelService {
	return &ChannelService{
		repo:        repo,
		twirbus:     twirbus,
		cfg:         cfg,
		cache:       cache,
		streamsrepo: streamsrepo,
	}
}

type ChannelService struct {
	repo        channelsrepo.Repository
	twirbus     *buscore.Bus
	cfg         config.Config
	cache       kv.KV
	streamsrepo streamsrepository.Repository
}

func createStreamsCacheKey(channelId uuid.UUID) string {
	return "twir:cache:channels:streams:" + channelId.String()
}

func (c *ChannelService) GetChannelByID(ctx context.Context, id uuid.UUID) (channelsmodel.Channel, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *ChannelService) GetChannelByConnectedUser(
	ctx context.Context,
	userID uuid.UUID,
	p platform.Platform,
) (channelsmodel.Channel, error) {
	switch p {
	case platform.PlatformTwitch:
		return c.repo.GetByTwitchUserID(ctx, userID)
	case platform.PlatformKick:
		return c.repo.GetByKickUserID(ctx, userID)
	default:
		return channelsmodel.Nil, fmt.Errorf("unknown platform: %s", p)
	}
}

func (c *ChannelService) GetChannelByPlatformUserID(
	ctx context.Context,
	platformUserID string,
	p platform.Platform,
) (channelsmodel.Channel, error) {
	switch p {
	case platform.PlatformTwitch:
		return c.repo.GetByTwitchPlatformID(ctx, platformUserID)
	case platform.PlatformKick:
		return c.repo.GetByKickPlatformID(ctx, platformUserID)
	default:
		return channelsmodel.Nil, fmt.Errorf("unknown platform: %s", p)
	}
}

type ChannelStreamWithChatLines struct {
	streamsmodel.Stream
	ParsedChatLines uint64
}

func (c *ChannelService) GetChannelStreams(
	ctx context.Context,
	channelID uuid.UUID,
) ([]ChannelStreamWithChatLines, error) {
	streams, err := c.streamsrepo.GetListByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	result := make([]ChannelStreamWithChatLines, len(streams))
	for i, s := range streams {
		result[i] = ChannelStreamWithChatLines{
			Stream: s,
		}
	}

	for i, stream := range streams {
		val, err := c.cache.Get(ctx, "stream:parsedMessages:"+stream.ID).Int()
		if err != nil {
			continue
		}

		result[i].ParsedChatLines = uint64(val)
	}

	return result, nil
}

func (c *ChannelService) IsChannelOnline(ctx context.Context, channelID uuid.UUID) (bool, error) {
	exists, err := c.cache.Exists(ctx, createStreamsCacheKey(channelID))
	if err != nil && !errors.Is(err, kv.ErrKeyNil) {
		return false, err
	}

	if err == nil {
		return exists, nil
	}

	streams, err := c.streamsrepo.GetListByChannelID(ctx, channelID)
	if err != nil {
		return false, err
	}

	online := len(streams) > 0
	if online {
		if err := c.cache.Set(
			ctx,
			createStreamsCacheKey(channelID),
			"1",
			kvoptions.WithExpire(30*time.Second),
		); err != nil {
			return false, err
		}
	}

	return online, nil
}

func (c *ChannelService) InvalidateOnlineCache(ctx context.Context, channelID uuid.UUID) error {
	return c.cache.Delete(ctx, createStreamsCacheKey(channelID))
}
