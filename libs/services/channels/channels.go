package channelservice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
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

func (c *ChannelService) GetChannelByID(ctx context.Context, id uuid.UUID) (channelentity.Channel, error) {
	return c.repo.GetByID(ctx, id)
}

// GetChannelByApiKey resolves a channel and its platform bindings from a channel API key.
func (c *ChannelService) GetChannelByApiKey(ctx context.Context, apiKey string) (channelentity.Channel, error) {
	return c.repo.GetByApiKey(ctx, apiKey)
}

// GetChannelByBindingUserID resolves a channel from a platform-scoped linked user ID.
func (c *ChannelService) GetChannelByBindingUserID(
	ctx context.Context,
	p platform.Platform,
	userID uuid.UUID,
) (channelentity.Channel, error) {
	return c.repo.GetByBindingUserID(ctx, p, userID)
}

// GetChannelByPlatformChannelID resolves a channel from a platform-scoped provider channel ID.
func (c *ChannelService) GetChannelByPlatformChannelID(
	ctx context.Context,
	p platform.Platform,
	platformChannelID string,
) (channelentity.Channel, error) {
	return c.repo.GetByPlatformChannelID(ctx, p, platformChannelID)
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

type GetBySlugOpts struct {
	Slug     string
	Platform *platform.Platform
}

func (c *ChannelService) GetBySlug(ctx context.Context, opts GetBySlugOpts) (channelentity.Channel, error) {
	channel, err := c.repo.GetBySlug(
		ctx,
		channelsrepo.GetBySlugInput{Slug: opts.Slug, Platform: opts.Platform},
	)
	if err != nil {
		return channelentity.Channel{}, err
	}

	return channel, nil
}
