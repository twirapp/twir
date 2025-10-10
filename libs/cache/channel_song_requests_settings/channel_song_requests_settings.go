package channelsongrequestssettings

import (
	"context"
	"time"

	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func New(
	db *gorm.DB,
	kv kv.KV,
) *generic_cacher.GenericCacher[model.ChannelSongRequestsSettings] {
	return generic_cacher.New[model.ChannelSongRequestsSettings](
		generic_cacher.Opts[model.ChannelSongRequestsSettings]{
			KV:        kv,
			KeyPrefix: "cache:twir:channelsongrequestssettings:channel:",
			LoadFn: func(ctx context.Context, key string) (model.ChannelSongRequestsSettings, error) {
				entity := &model.ChannelSongRequestsSettings{}
				err := db.
					WithContext(ctx).
					Where(`"channel_id" = ?`, key).
					Find(entity).
					Error
				if err != nil {
					return model.ChannelSongRequestsSettings{}, err
				}
				if entity.ID == "" {
					return model.ChannelSongRequestsSettings{}, nil
				}

				return *entity, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
