package giveaways

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channels_giveaways "github.com/twirapp/twir/libs/entities/channels_giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways"
)

func New(
	repo giveaways.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[[]channels_giveaways.Giveaway] {
	return generic_cacher.New(
		generic_cacher.Opts[[]channels_giveaways.Giveaway]{
			KV:        kv,
			KeyPrefix: "cache:twir:giveaways:channel:",
			LoadFn: func(ctx context.Context, key string) ([]channels_giveaways.Giveaway, error) {
				return repo.GetManyActiveByChannelID(ctx, key)
			},
			Ttl: 24 * 7 * time.Hour,
		},
	)
}
