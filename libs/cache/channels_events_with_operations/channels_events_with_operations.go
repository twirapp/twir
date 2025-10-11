package channelseventswithoperations

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/events"
	"github.com/twirapp/twir/libs/repositories/events/model"
)

func New(
	repo events.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[[]model.Event] {
	return generic_cacher.New(
		generic_cacher.Opts[[]model.Event]{
			KV:        kv,
			KeyPrefix: "cache:twir:channels_events_with_operations:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Event, error) {
				data, err := repo.GetManyByChannelID(ctx, key)
				if err != nil {
					return nil, err
				}

				return data, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
