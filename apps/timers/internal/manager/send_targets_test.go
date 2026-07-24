package manager

import (
	"encoding/json"
	"testing"

	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func TestGetTimerSendTargetsUsesBindingsByPlatform(t *testing.T) {
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{
			{
				Platform:          platformentity.PlatformKick,
				PlatformChannelID: "kick-channel",
				Enabled:           true,
			},
			{
				Platform:          platformentity.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				Enabled:           true,
				BotConfig:         json.RawMessage(`{"is_bot_mod":true}`),
			},
		},
	}

	want := []timerSendTarget{
		{platform: platformentity.PlatformTwitch, channelID: "twitch-channel"},
		{platform: platformentity.PlatformKick, channelID: "kick-channel"},
	}

	for _, timerPlatforms := range [][]platformentity.Platform{nil, {}} {
		assertTimerSendTargets(t, getTimerSendTargets(channel, timerPlatforms), want)
	}
}

func TestGetTimerSendTargetsFiltersBindingState(t *testing.T) {
	tests := []struct {
		name           string
		channel        channelentity.Channel
		timerPlatforms []platformentity.Platform
		want           []timerSendTarget
	}{
		{
			name: "restricts targets to configured platforms",
			channel: channelentity.Channel{
				Bindings: []channelplatformentity.ChannelPlatform{
					{
						Platform:          platformentity.PlatformTwitch,
						PlatformChannelID: "twitch-channel",
						Enabled:           true,
						BotConfig:         json.RawMessage(`{"is_bot_mod":true}`),
					},
					{
						Platform:          platformentity.PlatformKick,
						PlatformChannelID: "kick-channel",
						Enabled:           true,
					},
				},
			},
			timerPlatforms: []platformentity.Platform{platformentity.PlatformKick},
			want: []timerSendTarget{
				{platform: platformentity.PlatformKick, channelID: "kick-channel"},
			},
		},
		{
			name: "skips disabled bindings",
			channel: channelentity.Channel{
				Bindings: []channelplatformentity.ChannelPlatform{
					{
						Platform:          platformentity.PlatformTwitch,
						PlatformChannelID: "twitch-channel",
						Enabled:           false,
						BotConfig:         json.RawMessage(`{"is_bot_mod":true}`),
					},
					{
						Platform:          platformentity.PlatformKick,
						PlatformChannelID: "kick-channel",
						Enabled:           true,
					},
				},
			},
			want: []timerSendTarget{
				{platform: platformentity.PlatformKick, channelID: "kick-channel"},
			},
		},
		{
			name: "skips twitch binding without moderator state",
			channel: channelentity.Channel{
				Bindings: []channelplatformentity.ChannelPlatform{
					{
						Platform:          platformentity.PlatformTwitch,
						PlatformChannelID: "twitch-channel",
						Enabled:           true,
						BotConfig:         json.RawMessage(`{"is_bot_mod":false}`),
					},
				},
			},
			want: []timerSendTarget{},
		},
		{
			name: "skips malformed twitch configuration without blocking kick",
			channel: channelentity.Channel{
				Bindings: []channelplatformentity.ChannelPlatform{
					{
						Platform:          platformentity.PlatformTwitch,
						PlatformChannelID: "twitch-channel",
						Enabled:           true,
						BotConfig:         json.RawMessage(`{"is_bot_mod":`),
					},
					{
						Platform:          platformentity.PlatformKick,
						PlatformChannelID: "kick-channel",
						Enabled:           true,
					},
				},
			},
			want: []timerSendTarget{
				{platform: platformentity.PlatformKick, channelID: "kick-channel"},
			},
		},
		{
			name: "skips bindings without a provider channel id",
			channel: channelentity.Channel{
				Bindings: []channelplatformentity.ChannelPlatform{
					{
						Platform: platformentity.PlatformKick,
						Enabled:  true,
					},
				},
			},
			want: []timerSendTarget{},
		},
		{
			name: "skips unsupported bindings",
			channel: channelentity.Channel{
				Bindings: []channelplatformentity.ChannelPlatform{
					{
						Platform:          platformentity.PlatformVKVideoLive,
						PlatformChannelID: "vk-channel",
						Enabled:           true,
					},
				},
			},
			want: []timerSendTarget{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTimerSendTargets(t, getTimerSendTargets(tt.channel, tt.timerPlatforms), tt.want)
		})
	}
}

func TestHasSupportedTimerBinding(t *testing.T) {
	tests := []struct {
		name     string
		bindings []channelplatformentity.ChannelPlatform
		want     bool
	}{
		{
			name: "recognizes a disabled supported binding for initialization",
			bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platformentity.PlatformTwitch, PlatformChannelID: "twitch-channel"},
			},
			want: true,
		},
		{
			name: "recognizes a supported binding after an unsupported binding",
			bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platformentity.PlatformVKVideoLive, PlatformChannelID: "vk-channel", Enabled: true},
				{Platform: platformentity.PlatformKick, PlatformChannelID: "kick-channel"},
			},
			want: true,
		},
		{
			name: "ignores unsupported bindings",
			bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platformentity.PlatformVKVideoLive, PlatformChannelID: "vk-channel", Enabled: true},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channel := channelentity.Channel{Bindings: tt.bindings}
			if got := hasSupportedTimerBinding(channel); got != tt.want {
				t.Fatalf("hasSupportedTimerBinding() = %t, want %t", got, tt.want)
			}
		})
	}
}

func assertTimerSendTargets(t *testing.T, got, want []timerSendTarget) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d targets, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected target %+v at index %d, got %+v", want[i], i, got[i])
		}
	}
}
