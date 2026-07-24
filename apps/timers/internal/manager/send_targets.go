package manager

import (
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type timerSendTarget struct {
	platform  platformentity.Platform
	channelID string
}

func getTimerSendTargets(
	channel channelentity.Channel,
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

		binding, ok := channel.Binding(p)
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

func hasSupportedTimerBinding(channel channelentity.Channel) bool {
	for _, p := range []platformentity.Platform{
		platformentity.PlatformTwitch,
		platformentity.PlatformKick,
	} {
		if _, ok := channel.Binding(p); ok {
			return true
		}
	}

	return false
}

func isTwitchTimerBotMod(binding channelplatformentity.ChannelPlatform) bool {
	config, err := binding.ParseTwitchBotConfig()
	return err == nil && config.IsBotMod
}
