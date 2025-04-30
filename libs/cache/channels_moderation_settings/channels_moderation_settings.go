package channels_moderation_settings

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
)

func New(
	repo channels_moderation_settings.Repository,
	redis *redis.Client,
) *generic_cacher.GenericCacher[[]model.ChannelModerationSettings] {
	return generic_cacher.New[[]model.ChannelModerationSettings](
		generic_cacher.Opts[[]model.ChannelModerationSettings]{
			Redis:     redis,
			KeyPrefix: "cache:twir:channels_moderation_settings:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.ChannelModerationSettings, error) {
				return repo.GetByChannelID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
