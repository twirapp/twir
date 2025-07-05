package channelseventswithoperations

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	redis *redis.Client,
) *generic_cacher.GenericCacher[[]model.Event] {
	return generic_cacher.New[[]model.Event](
		generic_cacher.Opts[[]model.Event]{
			Redis:     redis,
			KeyPrefix: "cache:twir:channels_events_with_operations:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.Event, error) {
				var channelEvents []model.Event
				err := db.
					WithContext(ctx).
					Where(`"channelId" = ? AND "enabled" = ?`, key, true).
					Preload("Channel").
					Preload("Channel.User").
					Preload("Operations").
					Preload("Operations.Filters").
					Find(&channelEvents).
					Error
				if err != nil {
					return nil, err
				}

				return channelEvents, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
