package channelalerts

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/alerts"
	"github.com/twirapp/twir/libs/repositories/alerts/model"
)

func New(
	repository alerts.Repository,
	redis *redis.Client,
) *generic_cacher.GenericCacher[[]model.Alert] {
	return generic_cacher.New[[]model.Alert](
		generic_cacher.Opts[[]model.Alert]{
			Redis:     redis,
			KeyPrefix: "cache:twir:channels_alerts:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Alert, error) {
				return repository.GetManyByChannelID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
