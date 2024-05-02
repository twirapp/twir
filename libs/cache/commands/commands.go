package commands

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
) *db_generic_cacher.GenericCacher[[]model.ChannelsCommands] {
	return db_generic_cacher.New[[]model.ChannelsCommands](
		db_generic_cacher.Opts[[]model.ChannelsCommands]{
			Redis:     redis,
			KeyPrefix: "cache:twir:commands:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.ChannelsCommands, error) {
				var commands []model.ChannelsCommands
				err := db.
					WithContext(ctx).
					Model(&model.ChannelsCommands{}).
					Where(`channels_commands."channelId" = ? AND channels_commands."enabled" = ?`, key, true).
					Joins("Group").
					Preload("Responses").
					WithContext(ctx).
					Find(&commands).Error
				if err != nil {
					return nil, err
				}

				return commands, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
