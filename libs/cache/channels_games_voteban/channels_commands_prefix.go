package channels_games_voteban_cache

import (
	"context"
	"errors"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelsgamesvoteban "github.com/twirapp/twir/libs/repositories/channels_games_voteban"
	"github.com/twirapp/twir/libs/repositories/channels_games_voteban/model"
)

func New(
	repo channelsgamesvoteban.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[model.VoteBan] {
	return generic_cacher.New[model.VoteBan](
		generic_cacher.Opts[model.VoteBan]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:channels_games_voteban:channel:",
			LoadFn: func(ctx context.Context, key string) (model.VoteBan, error) {
				result, err := repo.GetByChannelID(ctx, key)
				if err != nil {
					if errors.Is(err, channelsgamesvoteban.ErrNotFound) {
						return model.Nil, nil
					}

					return model.Nil, err
				}

				return result, nil
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
