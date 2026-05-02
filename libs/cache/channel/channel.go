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
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
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
	usersRepo usersrepository.Repository,
	kv kv.KV,
) *TwitchUserIDCacher {
	return &TwitchUserIDCacher{
		GenericCacher: generic_cacher.New[channelmodel.Channel](
			generic_cacher.Opts[channelmodel.Channel]{
				KV:        kv,
				KeyPrefix: "cache:twir:channel_by_twitch_uid:",
				LoadFn: func(ctx context.Context, twitchUserID string) (channelmodel.Channel, error) {
					user, err := usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, twitchUserID)
					if err != nil {
						return channelmodel.Nil, fmt.Errorf("find user by twitch id %s: %w", twitchUserID, err)
					}
					if user.IsNil() {
						return channelmodel.Nil, usersmodel.ErrNotFound
					}

				return channelsRepo.GetByTwitchUserID(ctx, user.ID)
			},
				Ttl: 24 * time.Hour,
			},
		),
	}
}
