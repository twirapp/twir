package commands

import (
	"context"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

const KeyPrefix = "cache:twir:commands:channel:"

func New(
	db *gorm.DB,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[[]model.ChannelsCommands] {
	return generic_cacher.New[[]model.ChannelsCommands](
		generic_cacher.Opts[[]model.ChannelsCommands]{
			KV:        kvotter.New(),
			KeyPrefix: KeyPrefix,
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
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
