package handler

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	channelsmodel "github.com/twirapp/twir/libs/entities/channel"
	channelplatformsmodel "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestChannelBindingSelectsRequestedPlatformRegardlessOfOrder(t *testing.T) {
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platform.PlatformKick,
				UserID:            kickUserID,
				PlatformChannelID: "kick-channel",
			},
			{
				Platform:          platform.PlatformTwitch,
				UserID:            twitchUserID,
				PlatformChannelID: "twitch-channel",
			},
		},
	}

	binding, ok := channel.Binding(platform.PlatformTwitch)
	if !ok {
		t.Fatal("expected Twitch binding")
	}
	if binding.UserID != twitchUserID {
		t.Fatalf("selected user ID = %s, want %s", binding.UserID, twitchUserID)
	}
	if binding.PlatformChannelID != "twitch-channel" {
		t.Fatalf("selected platform channel ID = %q, want twitch-channel", binding.PlatformChannelID)
	}
}

func TestChannelPlatformParsesTwitchBotConfig(t *testing.T) {
	config, err := (channelplatformsmodel.ChannelPlatform{
		BotConfig: json.RawMessage(`{"bot_id":"bot-123","is_bot_mod":true,"is_twitch_banned":true}`),
	}).ParseTwitchBotConfig()
	if err != nil {
		t.Fatalf("ParseTwitchBotConfig returned error: %v", err)
	}
	if config.BotID != "bot-123" {
		t.Fatalf("bot ID = %q, want bot-123", config.BotID)
	}
	if !config.IsBotMod {
		t.Fatal("expected is_bot_mod to be true")
	}
	if !config.IsTwitchBanned {
		t.Fatal("expected is_twitch_banned to be true")
	}
}
