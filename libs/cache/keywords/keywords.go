package keywords

import (
	"context"
	"time"

	"github.com/google/uuid"
	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/repositories/keywords/model"

	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/keywords"
)

func New(
	repo keywords.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[[]model.Keyword] {
	return generic_cacher.New[[]model.Keyword](
		generic_cacher.Opts[[]model.Keyword]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:keywords:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Keyword, error) {
				parsedKey, err := uuid.Parse(key)
				if err != nil {
					return nil, err
				}

				return repo.GetAllByChannelID(ctx, parsedKey)
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
