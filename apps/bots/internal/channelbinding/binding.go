package channelbinding

import (
	"encoding/json"

	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type TwitchBotConfig struct {
	BotID          string `json:"bot_id"`
	IsBotMod       bool   `json:"is_bot_mod"`
	IsTwitchBanned bool   `json:"is_twitch_banned"`
}

func Find(
	channel channelsmodel.Channel,
	p platform.Platform,
) (channelplatformsmodel.ChannelPlatform, bool) {
	for _, binding := range channel.Bindings {
		if binding.Platform == p {
			return binding, true
		}
	}

	return channelplatformsmodel.ChannelPlatform{}, false
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

	if len(binding.BotConfig) == 0 {
		return binding, TwitchBotConfig{}, true, nil
	}

	var config TwitchBotConfig
	if err := json.Unmarshal(binding.BotConfig, &config); err != nil {
		return channelplatformsmodel.ChannelPlatform{}, TwitchBotConfig{}, false, err
	}

	return binding, config, true, nil
}
