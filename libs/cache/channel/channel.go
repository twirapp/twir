package channel

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func New(
	repo channelsrepository.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[channelmodel.Channel] {
	return generic_cacher.New[channelmodel.Channel](
		generic_cacher.Opts[channelmodel.Channel]{
			KV:        kv,
			KeyPrefix: "cache:twir:channel:",
			LoadFn: func(ctx context.Context, key string) (channelmodel.Channel, error) {
				return repo.GetByID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
