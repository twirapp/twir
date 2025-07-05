package channelsintegrationssettingsseventv

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	redis *redis.Client,
) *generic_cacher.GenericCacher[model.ChannelsIntegrationsSettingsSeventv] {
	return generic_cacher.New[model.ChannelsIntegrationsSettingsSeventv](
		generic_cacher.Opts[model.ChannelsIntegrationsSettingsSeventv]{
			Redis:     redis,
			KeyPrefix: "cache:twir:channelsintegrationssettingsseventv:channel:",
			LoadFn: func(ctx context.Context, key string) (
				model.ChannelsIntegrationsSettingsSeventv,
				error,
			) {
				settings := &model.ChannelsIntegrationsSettingsSeventv{}
				err := db.
					WithContext(ctx).
					Where(`"channel_id" = ?`, key).
					First(settings).
					Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return model.ChannelsIntegrationsSettingsSeventv{}, nil
					}
					return model.ChannelsIntegrationsSettingsSeventv{}, err
				}

				return *settings, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
