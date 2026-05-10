package manager

import (
	"testing"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestGetTimerSendTargets(t *testing.T) {
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	twitchPlatformID := "123"
	kickPlatformID := "kick-channel"

	baseChannel := channelmodel.Channel{
		TwitchUserID:     &twitchUserID,
		TwitchPlatformID: &twitchPlatformID,
		TwitchBotEnabled: true,
		KickUserID:       &kickUserID,
		KickPlatformID:   &kickPlatformID,
		KickBotEnabled:   true,
	}

	tests := []struct {
		name           string
		channel        channelmodel.Channel
		timerPlatforms []platformentity.Platform
		want           []timerSendTarget
	}{
		{
			name:           "empty platforms sends to all joined platforms",
			channel:        baseChannel,
			timerPlatforms: nil,
			want: []timerSendTarget{
				{platform: platformentity.PlatformTwitch, channelID: twitchPlatformID},
				{platform: platformentity.PlatformKick, channelID: kickPlatformID},
			},
		},
		{
			name:           "empty slice platforms sends to all joined platforms",
			channel:        baseChannel,
			timerPlatforms: []platformentity.Platform{},
			want: []timerSendTarget{
				{platform: platformentity.PlatformTwitch, channelID: twitchPlatformID},
				{platform: platformentity.PlatformKick, channelID: kickPlatformID},
			},
		},
		{
			name:           "explicit twitch only",
			channel:        baseChannel,
			timerPlatforms: []platformentity.Platform{platformentity.PlatformTwitch},
			want: []timerSendTarget{
				{platform: platformentity.PlatformTwitch, channelID: twitchPlatformID},
			},
		},
		{
			name:           "explicit kick only",
			channel:        baseChannel,
			timerPlatforms: []platformentity.Platform{platformentity.PlatformKick},
			want: []timerSendTarget{
				{platform: platformentity.PlatformKick, channelID: kickPlatformID},
			},
		},
		{
			name: "skips platforms where bot is not joined",
			channel: channelmodel.Channel{
				TwitchUserID:     &twitchUserID,
				TwitchPlatformID: &twitchPlatformID,
				TwitchBotEnabled: true,
				KickUserID:       &kickUserID,
				KickPlatformID:   &kickPlatformID,
				KickBotEnabled:   false,
			},
			timerPlatforms: nil,
			want: []timerSendTarget{
				{platform: platformentity.PlatformTwitch, channelID: twitchPlatformID},
			},
		},
		{
			name: "returns no targets when selected platform is not joined",
			channel: channelmodel.Channel{
				TwitchUserID:     &twitchUserID,
				TwitchPlatformID: &twitchPlatformID,
				TwitchBotEnabled: false,
				KickUserID:       &kickUserID,
				KickPlatformID:   &kickPlatformID,
				KickBotEnabled:   false,
			},
			timerPlatforms: []platformentity.Platform{platformentity.PlatformKick},
			want:           []timerSendTarget{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTimerSendTargets(tt.channel, tt.timerPlatforms)

			if len(got) != len(tt.want) {
				t.Fatalf("expected %d targets, got %d", len(tt.want), len(got))
			}

			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("expected target %+v at index %d, got %+v", tt.want[i], i, got[i])
				}
			}
		})
	}
}
