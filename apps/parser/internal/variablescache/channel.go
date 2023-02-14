package variables_cache

import (
	"github.com/samber/do"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
)

func (c *VariablesCacheService) GetTwitchChannel() *helix.ChannelInformation {
	cfg := do.MustInvoke[config.Config](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewAppClient(cfg, tokensGrpc)
	if err != nil {
		return nil
	}

	c.locks.twitchChannel.Lock()
	defer c.locks.twitchChannel.Unlock()

	if c.cache.TwitchChannel != nil {
		return c.cache.TwitchChannel
	}

	channel, err := twitchClient.GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: []string{c.ChannelId},
	})

	if err == nil && len(channel.Data.Channels) != 0 {
		c.cache.TwitchChannel = &channel.Data.Channels[0]
	}

	return c.cache.TwitchChannel
}
