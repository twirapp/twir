package commands

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Redis               *redis.Client
	CommandsPrefixCache *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
}

func New(opts Opts) *Service {
	return &Service{
		redis:               opts.Redis,
		commandsPrefixCache: opts.CommandsPrefixCache,
	}
}

type Service struct {
	redis               *redis.Client
	commandsPrefixCache *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
}

func (c *Service) GetCommandsPrefix(ctx context.Context, channelId string) (string, error) {
	var commandsPrefix string
	fetchedCommandsPrefix, err := c.commandsPrefixCache.Get(ctx, channelId)
	if err != nil && !errors.Is(err, channelscommandsprefixrepository.ErrNotFound) {
		return "", err
	}

	if fetchedCommandsPrefix == channelscommandsprefixmodel.Nil {
		commandsPrefix = "!"
	} else {
		commandsPrefix = fetchedCommandsPrefix.Prefix
	}

	return commandsPrefix, nil
}
