package manager

import (
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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

	if channel.TwitchBotJoined() &&
		channel.TwitchPlatformID != nil &&
		platformentity.ShouldExecute(timerPlatforms, platformentity.PlatformTwitch) {
		targets = append(targets, timerSendTarget{
			platform:  platformentity.PlatformTwitch,
			channelID: *channel.TwitchPlatformID,
		})
	}

	if channel.KickBotJoined() &&
		channel.KickPlatformID != nil &&
		platformentity.ShouldExecute(timerPlatforms, platformentity.PlatformKick) {
		targets = append(targets, timerSendTarget{
			platform:  platformentity.PlatformKick,
			channelID: *channel.KickPlatformID,
		})
	}

	return targets
}
