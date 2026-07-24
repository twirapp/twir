package events

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/events/internal/shared"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

type runtimeChannelRepositoryFake struct {
	channelsrepository.Repository

	channel                 channelentity.Channel
	lookupPlatform          platform.Platform
	lookupPlatformChannelID string
}

func (f *runtimeChannelRepositoryFake) GetByPlatformChannelID(
	_ context.Context,
	p platform.Platform,
	platformChannelID string,
) (channelentity.Channel, error) {
	f.lookupPlatform = p
	f.lookupPlatformChannelID = platformChannelID
	return f.channel, nil
}

func TestGetTwitchChannelRuntimeInfoSelectsTwitchBindingByPlatform(t *testing.T) {
	channelID := uuid.New()
	channel := channelentity.Channel{
		ID: channelID,
		Bindings: []channelplatformentity.ChannelPlatform{
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

func TestDualBoundKickEventKeepsEventIDAndResolvesTwitchRuntime(t *testing.T) {
	channelID := uuid.New()
	repo := &runtimeChannelRepositoryFake{
		channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{
				{
					Platform:          platform.PlatformKick,
					PlatformChannelID: "kick-channel",
				},
				{
					Platform:          platform.PlatformTwitch,
					PlatformChannelID: "twitch-channel",
					BotConfig:         json.RawMessage(`{"bot_id":"twitch-bot"}`),
				},
			},
		},
	}
	activity := Activity{
		channelService: channelservice.NewChannelService(
			repo,
			&buscore.Bus{},
			cfg.Config{},
			nil,
			nil,
		),
	}
	data := shared.EventData{
		ChannelID:               "kick-channel",
		ChannelTwitchPlatformID: "twitch-channel",
		Platform:                platform.PlatformKick,
	}

	runtimeChannel, err := activity.getTwitchChannelDbEntity(context.Background(), data)
	if err != nil {
		t.Fatalf("getTwitchChannelDbEntity returned error: %v", err)
	}
	if data.ChannelID != "kick-channel" {
		t.Errorf("event ChannelID = %q, want %q", data.ChannelID, "kick-channel")
	}
	if repo.lookupPlatform != platform.PlatformTwitch {
		t.Errorf("runtime lookup platform = %q, want %q", repo.lookupPlatform, platform.PlatformTwitch)
	}
	if repo.lookupPlatformChannelID != "twitch-channel" {
		t.Errorf("runtime lookup channel ID = %q, want %q", repo.lookupPlatformChannelID, "twitch-channel")
	}
	if runtimeChannel.ID != "twitch-channel" {
		t.Errorf("runtime broadcaster ID = %q, want %q", runtimeChannel.ID, "twitch-channel")
	}
}

func TestGetEventTwitchBotApiClientUsesTwitchBindingBotID(t *testing.T) {
	repo := &runtimeChannelRepositoryFake{
		channel: channelentity.Channel{
			ID: uuid.New(),
			Bindings: []channelplatformentity.ChannelPlatform{
				{
					Platform:          platform.PlatformKick,
					PlatformChannelID: "kick-channel",
				},
				{
					Platform:          platform.PlatformTwitch,
					PlatformChannelID: "twitch-channel",
					BotConfig:         json.RawMessage(`{"bot_id":"twitch-bot"}`),
				},
			},
		},
	}
	var clientBotID string
	activity := Activity{
		channelService: channelservice.NewChannelService(
			repo,
			&buscore.Bus{},
			cfg.Config{},
			nil,
			nil,
		),
		newTwitchBotClient: func(_ context.Context, botID string) (*helix.Client, error) {
			clientBotID = botID
			return &helix.Client{}, nil
		},
	}
	data := shared.EventData{
		ChannelID:               "kick-channel",
		ChannelTwitchPlatformID: "twitch-channel",
		Platform:                platform.PlatformKick,
	}

	client, err := activity.getEventTwitchBotApiClient(context.Background(), data)
	if err != nil {
		t.Fatalf("getEventTwitchBotApiClient returned error: %v", err)
	}
	if client == nil {
		t.Fatal("getEventTwitchBotApiClient returned nil client")
	}
	if clientBotID != "twitch-bot" {
		t.Errorf("bot client ID = %q, want %q", clientBotID, "twitch-bot")
	}
	if clientBotID == data.ChannelID {
		t.Errorf("bot client ID = event ChannelID %q, want Twitch bot ID", data.ChannelID)
	}
	if repo.lookupPlatform != platform.PlatformTwitch {
		t.Errorf("runtime lookup platform = %q, want %q", repo.lookupPlatform, platform.PlatformTwitch)
	}
	if repo.lookupPlatformChannelID != "twitch-channel" {
		t.Errorf("runtime lookup channel ID = %q, want %q", repo.lookupPlatformChannelID, "twitch-channel")
	}
}

func TestTwitchBroadcasterIDKeepsLegacyTwitchEventCompatibility(t *testing.T) {
	if got := twitchBroadcasterID(shared.EventData{
		ChannelID: "twitch-channel",
		Platform:  platform.PlatformTwitch,
	}); got != "twitch-channel" {
		t.Errorf("twitchBroadcasterID = %q, want %q", got, "twitch-channel")
	}
}

func TestTwitchBroadcasterIDDoesNotUsePlatformlessKickIDWithLegacyTwitchUser(t *testing.T) {
	if got := twitchBroadcasterID(shared.EventData{
		ChannelID:           "kick-channel",
		ChannelTwitchUserID: "twitch-user",
	}); got != "" {
		t.Errorf("twitchBroadcasterID = %q, want empty", got)
	}
}

func TestTwitchBroadcasterIDDoesNotUseKickEventID(t *testing.T) {
	if got := twitchBroadcasterID(shared.EventData{
		ChannelID: "kick-channel",
		Platform:  platform.PlatformKick,
	}); got != "" {
		t.Errorf("twitchBroadcasterID = %q, want empty", got)
	}
}

func TestTwitchBroadcasterIDDoesNotUsePlatformlessEventIDWithoutTwitchIdentity(t *testing.T) {
	if got := twitchBroadcasterID(shared.EventData{
		ChannelID: "unknown-channel",
	}); got != "" {
		t.Errorf("twitchBroadcasterID = %q, want empty", got)
	}
}
