package channelbinding

import (
	"encoding/json"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func Find(channel channelsmodel.Channel, p platform.Platform) (channelplatformsmodel.ChannelPlatform, bool) {
	for _, binding := range channel.Bindings {
		if binding.Platform == p {
			return binding, true
		}
	}

	return channelplatformsmodel.ChannelPlatform{}, false
}

type TwitchBotConfig struct {
	BotID          string `json:"bot_id"`
	IsBotMod       bool   `json:"is_bot_mod"`
	IsTwitchBanned bool   `json:"is_twitch_banned"`
}

func ParseTwitchBotConfig(binding channelplatformsmodel.ChannelPlatform) (TwitchBotConfig, error) {
	if len(binding.BotConfig) == 0 {
		return TwitchBotConfig{}, nil
	}

	var config TwitchBotConfig
	if err := json.Unmarshal(binding.BotConfig, &config); err != nil {
		return TwitchBotConfig{}, err
	}

	return config, nil
}

func FindTwitch(channel channelsmodel.Channel) (
	channelplatformsmodel.ChannelPlatform,
	TwitchBotConfig,
	bool,
	error,
) {
	binding, ok := Find(channel, platform.PlatformTwitch)
	if !ok {
		return channelplatformsmodel.ChannelPlatform{}, TwitchBotConfig{}, false, nil
	}

	config, err := ParseTwitchBotConfig(binding)
	if err != nil {
		return channelplatformsmodel.ChannelPlatform{}, TwitchBotConfig{}, false, err
	}

	return binding, config, true, nil
}

func NewParseContextChannel(
	channel channelsmodel.Channel,
	p platform.Platform,
	name string,
) (*types.ParseContextChannel, bool) {
	binding, ok := Find(channel, p)
	if !ok {
		return nil, false
	}

	parseContextChannel := &types.ParseContextChannel{
		ID:          binding.PlatformChannelID,
		Name:        name,
		DBChannelID: channel.ID.String(),
	}
	if twitchBinding, ok := Find(channel, platform.PlatformTwitch); ok {
		parseContextChannel.TwitchUserID = twitchBinding.UserID
	}

	return parseContextChannel, true
}
