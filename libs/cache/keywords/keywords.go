package keywords

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	"github.com/twirapp/twir/libs/repositories/keywords/model"

	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/keywords"
)

func New(
	repo keywords.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[[]model.Keyword] {
	return generic_cacher.New[[]model.Keyword](
		generic_cacher.Opts[[]model.Keyword]{
			KV:        kv,
			KeyPrefix: "cache:twir:keywords:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Keyword, error) {
				return repo.GetAllByChannelID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
