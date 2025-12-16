package channels_games_voteban_cache

import (
	"context"
	"errors"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/entities/voteban"
	channelsgamesvoteban "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
)

func New(
	repo channelsgamesvoteban.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[voteban.Voteban] {
	return generic_cacher.New[voteban.Voteban](
		generic_cacher.Opts[voteban.Voteban]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:channels_games_voteban:channel:",
			LoadFn: func(ctx context.Context, key string) (voteban.Voteban, error) {
				result, err := repo.GetByChannelID(ctx, key)
				if err != nil {
					if errors.Is(err, channelsgamesvoteban.ErrNotFound) {
						return voteban.Nil, nil
					}

					return voteban.Nil, err
				}

				return result, nil
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
