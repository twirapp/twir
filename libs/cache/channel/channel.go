package channel

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/kv"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func New(
	repo channelsrepository.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[channelmodel.Channel] {
	return generic_cacher.New[channelmodel.Channel](
		generic_cacher.Opts[channelmodel.Channel]{
			KV:        kv,
			KeyPrefix: "cache:twir:channel:",
			LoadFn: func(ctx context.Context, key string) (channelmodel.Channel, error) {
				parsed, err := uuid.Parse(key)
				if err != nil {
					return channelmodel.Nil, fmt.Errorf("invalid channel id: %w", err)
				}
				return repo.GetByID(ctx, parsed)
			},
			Ttl: 24 * time.Hour,
		},
	)
}

type TwitchUserIDCacher struct {
	*generic_cacher.GenericCacher[channelmodel.Channel]
}

func NewByTwitchUserID(
	channelsRepo channelsrepository.Repository,
	kv kv.KV,
) *TwitchUserIDCacher {
	return &TwitchUserIDCacher{
		GenericCacher: generic_cacher.New[channelmodel.Channel](
			generic_cacher.Opts[channelmodel.Channel]{
				KV:        kv,
				KeyPrefix: "cache:twir:channel_by_twitch_uid:",
				LoadFn: func(ctx context.Context, platformChannelID string) (channelmodel.Channel, error) {
					return channelsRepo.GetByPlatformChannelID(
						ctx,
						platform.PlatformTwitch,
						platformChannelID,
					)
				},
				Ttl: 24 * time.Hour,
			},
		),
	}
}
