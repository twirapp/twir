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

func TestNewParseContextChannelUsesSelectedPlatformAndTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	twitchUserID := uuid.New()
	channel := channelsmodel.Channel{
		ID: channelID,
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

	parseContextChannel, ok := NewParseContextChannel(channel, platform.PlatformKick, "channel-name", "")
	if !ok {
		t.Fatal("expected Kick parse context channel")
	}
	if parseContextChannel.ID != "kick-channel" {
		t.Errorf("ID = %q, want %q", parseContextChannel.ID, "kick-channel")
	}
	if parseContextChannel.DBChannelID != channelID.String() {
		t.Errorf("DBChannelID = %q, want %q", parseContextChannel.DBChannelID, channelID)
	}
	if parseContextChannel.TwitchUserID != twitchUserID {
		t.Errorf("TwitchUserID = %s, want %s", parseContextChannel.TwitchUserID, twitchUserID)
	}
}

func TestNewParseContextChannelUsesBindingIDWhenPresent(t *testing.T) {
	channelID := uuid.New()
	platformBinding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		Platform:          platform.PlatformKick,
		PlatformChannelID: "fallback-kick-channel",
	}
	canonicalBinding := channelplatformsmodel.ChannelPlatform{
		ID:                uuid.New(),
		Platform:          platform.PlatformKick,
		PlatformChannelID: "canonical-kick-channel",
	}
	channel := channelsmodel.Channel{
		ID:       channelID,
		Bindings: []channelplatformsmodel.ChannelPlatform{platformBinding, canonicalBinding},
	}

	parseContextChannel, ok := NewParseContextChannel(
		channel,
		platform.PlatformKick,
		"channel-name",
		canonicalBinding.ID.String(),
	)
	if !ok {
		t.Fatal("expected canonical Kick parse context channel")
	}
	if parseContextChannel.ID != canonicalBinding.PlatformChannelID {
		t.Errorf("ID = %q, want %q", parseContextChannel.ID, canonicalBinding.PlatformChannelID)
	}

	_, ok = NewParseContextChannel(channel, platform.PlatformKick, "channel-name", uuid.New().String())
	if ok {
		t.Fatal("expected no parse context channel for an unknown binding ID")
	}

	_, ok = NewParseContextChannel(channel, platform.PlatformKick, "channel-name", "not-a-uuid")
	if ok {
		t.Fatal("expected no parse context channel for an invalid binding ID")
	}
}

func TestFindTwitchUsesTwitchBotConfig(t *testing.T) {
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:  platform.PlatformKick,
				BotConfig: json.RawMessage(`{"bot_id":"kick-bot"}`),
			},
			{
				Platform: platform.PlatformTwitch,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
		},
	}

	binding, config, ok, err := FindTwitch(channel)
	if err != nil {
		t.Fatalf("FindTwitch returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected Twitch binding")
	}
	if binding.Platform != platform.PlatformTwitch {
		t.Errorf("Platform = %q, want %q", binding.Platform, platform.PlatformTwitch)
	}
	if config.BotID != "twitch-bot" {
		t.Errorf("BotID = %q, want %q", config.BotID, "twitch-bot")
	}
	if !config.IsBotMod {
		t.Error("IsBotMod = false, want true")
	}
	if !config.IsTwitchBanned {
		t.Error("IsTwitchBanned = false, want true")
	}
}
