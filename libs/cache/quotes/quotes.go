package quotes

import (
	"context"
	"time"

	"github.com/google/uuid"
	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/quotes"
	"github.com/twirapp/twir/libs/repositories/quotes/model"
)

func New(
	repo quotes.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[[]model.Quote] {
	return generic_cacher.New(
		generic_cacher.Opts[[]model.Quote]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:quotes:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Quote, error) {
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
