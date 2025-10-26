package channels_commands_prefix

import (
	"context"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
)

func New(
	repo channels_commands_prefix.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[model.ChannelsCommandsPrefix] {
	return generic_cacher.New[model.ChannelsCommandsPrefix](
		generic_cacher.Opts[model.ChannelsCommandsPrefix]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:channels_commands_prefix:channel:",
			LoadFn: func(ctx context.Context, key string) (model.ChannelsCommandsPrefix, error) {
				return repo.GetByChannelID(ctx, key)
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
