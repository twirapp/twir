package manager

import (
	"encoding/json"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type timerSendTarget struct {
	platform  platformentity.Platform
	channelID string
}

func getTimerSendTargets(
	channel channelmodel.Channel,
	timerPlatforms []platformentity.Platform,
) []timerSendTarget {
	targets := make([]timerSendTarget, 0, 2)

	for _, p := range []platformentity.Platform{
		platformentity.PlatformTwitch,
		platformentity.PlatformKick,
	} {
		if !platformentity.ShouldExecute(timerPlatforms, p) {
			continue
		}

		binding, ok := findTimerBinding(channel, p)
		if !ok || !binding.Enabled || binding.PlatformChannelID == "" {
			continue
		}

		if p == platformentity.PlatformTwitch && !isTwitchTimerBotMod(binding) {
			continue
		}

		targets = append(targets, timerSendTarget{
			platform:  p,
			channelID: binding.PlatformChannelID,
		})
	}

	return targets
}

func hasSupportedTimerBinding(channel channelmodel.Channel) bool {
	for _, p := range []platformentity.Platform{
		platformentity.PlatformTwitch,
		platformentity.PlatformKick,
	} {
		if _, ok := findTimerBinding(channel, p); ok {
			return true
		}
	}

	return false
}

func findTimerBinding(
	channel channelmodel.Channel,
	p platformentity.Platform,
) (channelplatformsmodel.ChannelPlatform, bool) {
	for _, binding := range channel.Bindings {
		if binding.Platform == p {
			return binding, true
		}
	}

	return channelplatformsmodel.ChannelPlatform{}, false
}

func isTwitchTimerBotMod(binding channelplatformsmodel.ChannelPlatform) bool {
	var config struct {
		IsBotMod bool `json:"is_bot_mod"`
	}

	return json.Unmarshal(binding.BotConfig, &config) == nil && config.IsBotMod
}
