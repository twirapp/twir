package events

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestGetTwitchChannelRuntimeInfoSelectsTwitchBindingByPlatform(t *testing.T) {
	channelID := uuid.New()
	channel := channelsmodel.Channel{
		ID: channelID,
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platform.PlatformKick,
				PlatformChannelID: "kick-channel",
				Enabled:           true,
			},
			{
				Platform:          platform.PlatformTwitch,
				PlatformChannelID: "twitch-channel",
				Enabled:           true,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
		},
	}

	info, err := getTwitchChannelRuntimeInfo(channel)
	if err != nil {
		t.Fatalf("getTwitchChannelRuntimeInfo returned error: %v", err)
	}
	if info.ChannelID != channelID.String() {
		t.Errorf("ChannelID = %q, want %q", info.ChannelID, channelID)
	}
	if info.BroadcasterUserID != "twitch-channel" {
		t.Errorf("BroadcasterUserID = %q, want %q", info.BroadcasterUserID, "twitch-channel")
	}
	if info.TwitchPlatformID != "twitch-channel" {
		t.Errorf("TwitchPlatformID = %q, want %q", info.TwitchPlatformID, "twitch-channel")
	}
	if info.BotID != "twitch-bot" {
		t.Errorf("BotID = %q, want %q", info.BotID, "twitch-bot")
	}
	if !info.IsBotMod {
		t.Error("IsBotMod = false, want true")
	}
	if !info.IsTwitchBanned {
		t.Error("IsTwitchBanned = false, want true")
	}
}
