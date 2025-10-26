package giveaways

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways/model"
)

func New(
	repo giveaways.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[[]model.ChannelGiveaway] {
	return generic_cacher.New(
		generic_cacher.Opts[[]model.ChannelGiveaway]{
			KV:        kv,
			KeyPrefix: "cache:twir:giveaways:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.ChannelGiveaway, error) {
				return repo.GetManyActiveByChannelID(ctx, key)
			},
			Ttl: 24 * 7 * time.Hour,
		},
	)
}
