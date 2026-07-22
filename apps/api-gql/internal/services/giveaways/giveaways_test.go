package giveaways

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/channels_giveaways_settings"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestSendWinnerMessageUsesSelectedTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	service := &Service{
		giveawaysSettingsRepository: giveawayTestSettingsRepository{
			settings: channels_giveaways_settings.Settings{WinnerMessage: "winner"},
		},
		channelService: giveawayTestChannelLookup{channel: channelsmodel.Channel{
			ID: channelID,
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{
					Platform:  platform.PlatformVKVideoLive,
					BotConfig: json.RawMessage(`{"unexpected":"vk"}`),
				},
				{
					Platform:  platform.PlatformKick,
					BotConfig: json.RawMessage(`{"unexpected":"kick"}`),
				},
				{
					Platform:          platform.PlatformTwitch,
					PlatformChannelID: "",
					BotConfig:         json.RawMessage(`{"unexpected":"twitch"}`),
				},
			},
		}},
	}

	err := service.sendWinnerMessage(context.Background(), channelID.String(), nil)
	if err != nil {
		t.Fatalf("send winner message: %v", err)
	}
}

func TestSendWinnerMessageRejectsMissingTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	service := &Service{
		giveawaysSettingsRepository: giveawayTestSettingsRepository{
			settings: channels_giveaways_settings.Settings{WinnerMessage: "winner"},
		},
		channelService: giveawayTestChannelLookup{channel: channelsmodel.Channel{
			ID: channelID,
			Bindings: []channelplatformsmodel.ChannelPlatform{
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
	channel channelsmodel.Channel
}

func (r giveawayTestChannelLookup) GetChannelByID(
	context.Context,
	uuid.UUID,
) (channelsmodel.Channel, error) {
	return r.channel, nil
}
