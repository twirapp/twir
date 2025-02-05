package greetings

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/greetings"
	"github.com/twirapp/twir/libs/repositories/greetings/model"
)

func New(
	repo greetings.Repository,
	redis *redis.Client,
) *generic_cacher.GenericCacher[[]model.Greeting] {
	return generic_cacher.New[[]model.Greeting](
		generic_cacher.Opts[[]model.Greeting]{
			Redis:     redis,
			KeyPrefix: "cache:twir:greetings:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Greeting, error) {
				return repo.GetManyByChannelID(ctx, key, greetings.GetManyInput{})
			},
			Ttl: 24 * time.Hour,
		},
	)
}
