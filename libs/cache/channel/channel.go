package channel

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func New(
	repo channelsrepository.Repository,
	redis *redis.Client,
) *generic_cacher.GenericCacher[channelmodel.Channel] {
	return generic_cacher.New[channelmodel.Channel](
		generic_cacher.Opts[channelmodel.Channel]{
			Redis:     redis,
			KeyPrefix: "cache:twir:channel:",
			LoadFn: func(ctx context.Context, key string) (channelmodel.Channel, error) {
				return repo.GetByID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
