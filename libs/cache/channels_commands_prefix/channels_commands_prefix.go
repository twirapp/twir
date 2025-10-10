package channels_commands_prefix

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
)

func New(
	repo channels_commands_prefix.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[model.ChannelsCommandsPrefix] {
	return generic_cacher.New[model.ChannelsCommandsPrefix](
		generic_cacher.Opts[model.ChannelsCommandsPrefix]{
			KV:        kv,
			KeyPrefix: "cache:twir:channels_commands_prefix:channel:",
			LoadFn: func(ctx context.Context, key string) (model.ChannelsCommandsPrefix, error) {
				return repo.GetByChannelID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
