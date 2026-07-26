package workflows

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/events/internal/shared"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestGetEventChannelBindingsSelectsEventAndTwitchBindingsByPlatform(t *testing.T) {
	twitchUserID := uuid.New()
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				UserID:            twitchUserID,
				Enabled:           false,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
				Enabled:           true,
			},
		},
	}

	bindings, err := getEventChannelBindings(channel, platform.PlatformKick)
	if err != nil {
		t.Fatalf("getEventChannelBindings returned error: %v", err)
	}
	if bindings.event.PlatformChannelID != "kick-channel" {
		t.Errorf("event PlatformChannelID = %q, want %q", bindings.event.PlatformChannelID, "kick-channel")
	}
	if !bindings.event.Enabled {
		t.Error("event binding Enabled = false, want true")
	}
	if !bindings.hasTwitch {
		t.Fatal("expected Twitch binding")
	}
	if bindings.twitch.UserID != twitchUserID {
		t.Errorf("Twitch UserID = %s, want %s", bindings.twitch.UserID, twitchUserID)
	}
	if bindings.twitchBotConfig.BotID != "twitch-bot" {
		t.Errorf("Twitch bot ID = %q, want %q", bindings.twitchBotConfig.BotID, "twitch-bot")
	}
	if !bindings.twitchBotConfig.IsBotMod {
		t.Error("Twitch bot mod state = false, want true")
	}
	if !bindings.twitchBotConfig.IsTwitchBanned {
		t.Error("Twitch ban state = false, want true")
	}

	data := bindings.applyTo(shared.EventData{
		ChannelID: "incoming-kick-channel",
		Platform:  platform.PlatformKick,
	})
	if data.ChannelID != "kick-channel" {
		t.Errorf("event ChannelID = %q, want %q", data.ChannelID, "kick-channel")
	}
	if data.ChannelTwitchPlatformID != "twitch-channel" {
		t.Errorf("Twitch broadcaster ID = %q, want %q", data.ChannelTwitchPlatformID, "twitch-channel")
	}
	if data.ChannelTwitchUserID != twitchUserID.String() {
		t.Errorf("Twitch user ID = %q, want %q", data.ChannelTwitchUserID, twitchUserID)
	}
}

func TestIsTwitchBannedEventScopesBanStateToTwitch(t *testing.T) {
	tests := []struct {
		name           string
		eventPlatform  platform.Platform
		wantSuppressed bool
	}{
		{
			name:           "Twitch event is suppressed",
			eventPlatform:  platform.PlatformTwitch,
			wantSuppressed: true,
		},
		{
			name:           "Kick event flows",
			eventPlatform:  platform.PlatformKick,
			wantSuppressed: false,
		},
		{
			name:           "VK event flows",
			eventPlatform:  platform.PlatformVKVideoLive,
			wantSuppressed: false,
		},
		{
			name:           "platform-neutral event flows",
			eventPlatform:  "",
			wantSuppressed: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isTwitchBanned := true
			if got := isTwitchBannedEvent(test.eventPlatform, isTwitchBanned); got != test.wantSuppressed {
				t.Errorf("isTwitchBannedEvent(%q, true) = %t, want %t", test.eventPlatform, got, test.wantSuppressed)
			}
		})
	}
}
