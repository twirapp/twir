package channels_commands_prefix

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/cache"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
)

func NewInMemory(
	repo channels_commands_prefix.Repository,
) (cache.InMemory[model.ChannelsCommandsPrefix], error) {
	return cache.NewInMemory(
		cache.InMemoryOptions{
			MaxSize:     5000,
			MinCapacity: 1000,
		},
		15*time.Minute,
		func(ctx context.Context, channelId cache.Key) (model.ChannelsCommandsPrefix, error) {
			prefix, err := repo.GetByChannelID(ctx, channelId)
			if err != nil {
				if !errors.Is(err, channels_commands_prefix.ErrNotFound) {
					return model.ChannelsCommandsPrefix{}, err
				}

				// Create channel command prefix if it's not exists yet.
				prefix, err = repo.Create(
					ctx, channels_commands_prefix.CreateInput{
						ChannelID: channelId,
						Prefix:    "!",
					},
				)
				if err != nil && !errors.Is(err, channels_commands_prefix.ErrAlreadyExists) {
					return model.ChannelsCommandsPrefix{}, err
				}
			}

			return prefix, nil
		},
	)
}

func New(
	repo channels_commands_prefix.Repository,
	redis *redis.Client,
) *generic_cacher.GenericCacher[model.ChannelsCommandsPrefix] {
	return generic_cacher.New[model.ChannelsCommandsPrefix](
		generic_cacher.Opts[model.ChannelsCommandsPrefix]{
			Redis:     redis,
			KeyPrefix: "cache:twir:channels_commands_prefix:channel:",
			LoadFn: func(ctx context.Context, key string) (model.ChannelsCommandsPrefix, error) {
				return repo.GetByChannelID(ctx, key)
			},
			Ttl: 24 * time.Hour,
		},
	)
}
