package chat_wall

import (
	"context"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/chat_wall"
	"github.com/twirapp/twir/libs/repositories/chat_wall/model"
)

func NewEnabledOnly(
	repo chat_wall.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[[]model.ChatWall] {
	return generic_cacher.New[[]model.ChatWall](
		generic_cacher.Opts[[]model.ChatWall]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:channels_chat_wall:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.ChatWall, error) {
				enabled := true

				return repo.GetMany(
					ctx,
					chat_wall.GetManyInput{
						ChannelID: key,
						Enabled:   &enabled,
					},
				)
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}

func NewSettings(
	repo chat_wall.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[model.ChatWallSettings] {
	return generic_cacher.New[model.ChatWallSettings](
		generic_cacher.Opts[model.ChatWallSettings]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:channels_chat_wall_settings:channel:",
			LoadFn: func(ctx context.Context, key string) (model.ChatWallSettings, error) {
				return repo.GetChannelSettings(
					ctx,
					key,
				)
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
