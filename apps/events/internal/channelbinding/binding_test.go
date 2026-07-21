package channelbinding

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestFindSelectsBindingByPlatform(t *testing.T) {
	twitchUserID := uuid.New()
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
			},
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				UserID:            twitchUserID,
			},
		},
	}

	binding, ok := Find(channel, platform.PlatformTwitch)
	if !ok {
		t.Fatal("expected Twitch binding")
	}
	if binding.PlatformChannelID != "twitch-channel" {
		t.Errorf("PlatformChannelID = %q, want %q", binding.PlatformChannelID, "twitch-channel")
	}
	if binding.UserID != twitchUserID {
		t.Errorf("UserID = %s, want %s", binding.UserID, twitchUserID)
	}
}

func TestParseTwitchBotConfig(t *testing.T) {
	config, err := ParseTwitchBotConfig(channelplatformsmodel.ChannelPlatform{
		BotConfig: json.RawMessage(`{"bot_id":"bot-123","is_bot_mod":true,"is_twitch_banned":true}`),
	})
	if err != nil {
		t.Fatalf("ParseTwitchBotConfig returned error: %v", err)
	}
	if config.BotID != "bot-123" {
		t.Errorf("BotID = %q, want %q", config.BotID, "bot-123")
	}
	if !config.IsBotMod {
		t.Error("IsBotMod = false, want true")
	}
	if !config.IsTwitchBanned {
		t.Error("IsTwitchBanned = false, want true")
	}
}

func TestParseTwitchBotConfigEmpty(t *testing.T) {
	config, err := ParseTwitchBotConfig(channelplatformsmodel.ChannelPlatform{})
	if err != nil {
		t.Fatalf("ParseTwitchBotConfig returned error: %v", err)
	}
	if config != (TwitchBotConfig{}) {
		t.Errorf("config = %+v, want zero value", config)
	}
}
