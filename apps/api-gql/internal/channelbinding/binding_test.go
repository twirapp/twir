package channelbinding

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestFindTwitchSelectsTwitchBindingRegardlessOfOrder(t *testing.T) {
	twitchUserID := uuid.New()
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
				Enabled:           false,
			},
			{
				Platform:          platform.PlatformTwitch,
				UserID:            twitchUserID,
				PlatformChannelID: "twitch-channel",
				Enabled:           true,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
		},
	}

	binding, config, found, err := FindTwitch(channel)
	if err != nil {
		t.Fatalf("find Twitch binding: %v", err)
	}
	if !found {
		t.Fatal("expected Twitch binding")
	}
	if binding.PlatformChannelID != "twitch-channel" {
		t.Fatalf("platform channel ID = %q, want Twitch binding", binding.PlatformChannelID)
	}
	if binding.UserID != twitchUserID {
		t.Fatalf("user ID = %s, want %s", binding.UserID, twitchUserID)
	}
	if config.BotID != "twitch-bot" || !config.IsBotMod || !config.IsTwitchBanned {
		t.Fatalf("Twitch bot config = %#v, want mapped Twitch config", config)
	}
}

func TestFindTwitchRejectsMalformedBotConfig(t *testing.T) {
	_, _, _, err := FindTwitch(channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{{
			Platform:  platform.PlatformTwitch,
			BotConfig: json.RawMessage(`{"bot_id":`),
		}},
	})
	if err == nil {
		t.Fatal("expected malformed Twitch bot config error")
	}
}
