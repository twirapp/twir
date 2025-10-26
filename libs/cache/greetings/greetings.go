package greetings

import (
	"context"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/greetings"
	"github.com/twirapp/twir/libs/repositories/greetings/model"
)

func New(
	repo greetings.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[[]model.Greeting] {
	return generic_cacher.New[[]model.Greeting](
		generic_cacher.Opts[[]model.Greeting]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:greetings:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Greeting, error) {
				return repo.GetManyByChannelID(ctx, key, greetings.GetManyInput{})
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
