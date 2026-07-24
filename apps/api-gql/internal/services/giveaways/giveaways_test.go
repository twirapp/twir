package giveaways

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	botsbus "github.com/twirapp/twir/libs/bus-core/bots"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/channels_giveaways"
	"github.com/twirapp/twir/libs/entities/channels_giveaways_settings"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestSendWinnerMessagePublishesSelectedTwitchNATSRequest(t *testing.T) {
	channelID := uuid.New()
	twitchBinding := channelplatformentity.ChannelPlatform{
		ChannelID:         channelID,
		Platform:          platform.PlatformTwitch,
		PlatformChannelID: "",
		BotConfig:         json.RawMessage(`{"unexpected":"twitch"}`),
	}
	publisher := &giveawayTestWinnerMessagePublisher{}
	service := &Service{
		giveawaysSettingsRepository: giveawayTestSettingsRepository{
			settings: channels_giveaways_settings.Settings{WinnerMessage: "winner {winner}"},
		},
		channelService: giveawayTestChannelLookup{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{
				{
					ChannelID:         uuid.New(),
					Platform:          platform.PlatformVKVideoLive,
					PlatformChannelID: "wrong-vk-channel",
					BotConfig:         json.RawMessage(`{"unexpected":"vk"}`),
				},
				{
					ChannelID:         uuid.New(),
					Platform:          platform.PlatformKick,
					PlatformChannelID: "wrong-kick-channel",
					BotConfig:         json.RawMessage(`{"unexpected":"kick"}`),
				},
				twitchBinding,
			},
		}},
		winnerMessagePublisher: publisher,
	}

	err := service.sendWinnerMessage(context.Background(), channelID.String(), []channels_giveaways.GiveawayWinner{
		{UserID: uuid.New(), UserLogin: "winner"},
	})
	if err != nil {
		t.Fatalf("send winner message: %v", err)
	}
	if len(publisher.requests) != 1 {
		t.Fatalf("published requests = %d, want 1", len(publisher.requests))
	}

	request := publisher.requests[0]
	if request.ChannelID != twitchBinding.ChannelID {
		t.Fatalf("request channel ID = %s, want selected Twitch channel ID %s", request.ChannelID, twitchBinding.ChannelID)
	}
	if len(request.Platforms) != 1 || request.Platforms[0] != platform.PlatformTwitch {
		t.Fatalf("request platforms = %#v, want only Twitch", request.Platforms)
	}
	if !request.SkipRateLimits {
		t.Fatal("expected winner message to skip rate limits")
	}
	if request.Message != "winner winner" {
		t.Fatalf("request message = %q, want winner replacement", request.Message)
	}
}

func TestSendWinnerMessageRejectsMissingTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	service := &Service{
		giveawaysSettingsRepository: giveawayTestSettingsRepository{
			settings: channels_giveaways_settings.Settings{WinnerMessage: "winner"},
		},
		channelService: giveawayTestChannelLookup{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{
				{Platform: platform.PlatformVKVideoLive},
				{Platform: platform.PlatformKick},
			},
		}},
	}

	if err := service.sendWinnerMessage(context.Background(), channelID.String(), nil); err == nil {
		t.Fatal("expected missing Twitch binding error")
	}
}

type giveawayTestSettingsRepository struct {
	settings channels_giveaways_settings.Settings
}

func (r giveawayTestSettingsRepository) GetByChannelID(
	context.Context,
	string,
) (channels_giveaways_settings.Settings, error) {
	return r.settings, nil
}

func (giveawayTestSettingsRepository) Update(
	context.Context,
	string,
	channels_giveaways_settings.Settings,
) (channels_giveaways_settings.Settings, error) {
	return channels_giveaways_settings.Settings{}, nil
}

type giveawayTestChannelLookup struct {
	channel channelentity.Channel
}

func (r giveawayTestChannelLookup) GetChannelByID(
	context.Context,
	uuid.UUID,
) (channelentity.Channel, error) {
	return r.channel, nil
}

type giveawayTestWinnerMessagePublisher struct {
	requests []botsbus.SendMessageRequest
}

func (p *giveawayTestWinnerMessagePublisher) Publish(
	_ context.Context,
	request botsbus.SendMessageRequest,
) error {
	p.requests = append(p.requests, request)
	return nil
}
