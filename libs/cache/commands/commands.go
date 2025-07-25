package commands

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	model "github.com/twirapp/twir/libs/gomodels"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	redis *redis.Client,
) *generic_cacher.GenericCacher[[]model.ChannelsCommands] {
	return generic_cacher.New[[]model.ChannelsCommands](
		generic_cacher.Opts[[]model.ChannelsCommands]{
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
