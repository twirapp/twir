package channels_moderation_settings

import (
	"context"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings"
	"github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
)

func New(
	repo channels_moderation_settings.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[[]model.ChannelModerationSettings] {
	return generic_cacher.New[[]model.ChannelModerationSettings](
		generic_cacher.Opts[[]model.ChannelModerationSettings]{
			KV:        kvotter.New(),
			KeyPrefix: "cache:twir:channels_moderation_settings:channel:",
			LoadFn: func(ctx context.Context, key string) ([]model.ChannelModerationSettings, error) {
				return repo.GetByChannelID(ctx, key)
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
