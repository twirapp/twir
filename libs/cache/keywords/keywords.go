package keywords

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	db_generic_cacher "github.com/twirapp/twir/libs/cache/db-generic-cacher"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	redis *redis.Client,
) *db_generic_cacher.GenericCacher[[]model.ChannelsKeywords] {
	return db_generic_cacher.New[[]model.ChannelsKeywords](
		db_generic_cacher.Opts[[]model.ChannelsKeywords]{
			Redis:     redis,
			KeyPrefix: "cache:twir:keywords:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.ChannelsKeywords, error) {
				var keywords []model.ChannelsKeywords
				err := db.WithContext(ctx).Where(
					`"channelId" = ? AND "enabled" = ?`, key,
					true,
				).Find(&keywords).Error
				if err != nil {
					return nil, err
				}
				return keywords, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
