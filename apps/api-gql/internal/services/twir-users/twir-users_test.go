package twir_users

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"github.com/twirapp/twir/libs/repositories/users_with_channel"
	"github.com/twirapp/twir/libs/repositories/users_with_channel/model"
)

type usersWithChannelsRepositoryStub struct {
	getManyCalls int
	countCalls   int
}

func (s *usersWithChannelsRepositoryStub) GetByID(context.Context, string) (model.UserWithChannel, error) {
	return model.UserWithChannel{}, nil
}

func (s *usersWithChannelsRepositoryStub) GetManyByIDS(
	context.Context,
	users_with_channel.GetManyInput,
) ([]model.UserWithChannel, error) {
	s.getManyCalls++
	return nil, nil
}

func (s *usersWithChannelsRepositoryStub) GetManyCount(context.Context, users_with_channel.GetManyInput) (int, error) {
	s.countCalls++
	return 0, nil
}

func TestGetManyRejectsInvalidBadgeID(t *testing.T) {
	repository := &usersWithChannelsRepositoryStub{}
	service := &Service{usersWithChannelsRepository: repository}

	_, err := service.GetMany(context.Background(), GetManyInput{
		HasBadges: []string{"not-a-uuid"},
	})
	if err == nil {
		t.Fatal("expected invalid badge ID to be rejected")
	}
	if repository.getManyCalls != 0 || repository.countCalls != 0 {
		t.Fatalf(
			"expected no repository calls, got getMany=%d count=%d",
			repository.getManyCalls,
			repository.countCalls,
		)
	}
}

func TestModelToEntityUsesMatchingTwitchBinding(t *testing.T) {
	channelID := uuid.New()
	kickBotUserID := uuid.New()
	result, err := (&Service{}).modelToEntity(model.UserWithChannel{
		User: usersmodel.User{ID: uuid.New(), Platform: platform.PlatformTwitch},
		Channel: &channelsmodel.Channel{
			ID: channelID,
			Bindings: []channelplatformsmodel.ChannelPlatform{
				{
					Platform:          platform.PlatformKick,
					PlatformChannelID: "kick-channel",
					Enabled:           false,
					BotUserID:         &kickBotUserID,
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
		},
	})
	if err != nil {
		t.Fatalf("map user with channel: %v", err)
	}
	if result.Channel == nil {
		t.Fatal("expected channel")
	}
	if result.Channel.ID != channelID {
		t.Fatalf("channel ID = %s, want %s", result.Channel.ID, channelID)
	}
	if !result.Channel.IsEnabled {
		t.Fatal("expected Twitch binding enabled state")
	}
	if !result.Channel.IsBotMod || !result.Channel.IsTwitchBanned {
		t.Fatalf("Twitch state = %#v, want parsed Twitch bot config", result.Channel)
	}
	if result.Channel.BotID != "twitch-bot" {
		t.Fatalf("bot ID = %q, want Twitch bot ID", result.Channel.BotID)
	}
}

func TestModelToEntityUsesBindingBotUserIDForNonTwitchPlatform(t *testing.T) {
	botUserID := uuid.New()
	result, err := (&Service{}).modelToEntity(model.UserWithChannel{
		User: usersmodel.User{ID: uuid.New(), Platform: platform.PlatformKick},
		Channel: &channelsmodel.Channel{
			ID: uuid.New(),
			Bindings: []channelplatformsmodel.ChannelPlatform{{
				Platform:  platform.PlatformKick,
				Enabled:   true,
				BotUserID: &botUserID,
			}},
		},
	})
	if err != nil {
		t.Fatalf("map user with channel: %v", err)
	}
	if result.Channel == nil {
		t.Fatal("expected channel")
	}
	if !result.Channel.IsEnabled {
		t.Fatal("expected Kick binding enabled state")
	}
	if result.Channel.BotID != botUserID.String() {
		t.Fatalf("bot ID = %q, want %s", result.Channel.BotID, botUserID)
	}
}

func TestModelToEntityRejectsMalformedTwitchBotConfig(t *testing.T) {
	_, err := (&Service{}).modelToEntity(model.UserWithChannel{
		User: usersmodel.User{ID: uuid.New(), Platform: platform.PlatformTwitch},
		Channel: &channelsmodel.Channel{
			ID: uuid.New(),
			Bindings: []channelplatformsmodel.ChannelPlatform{{
				Platform:  platform.PlatformTwitch,
				BotConfig: json.RawMessage(`{"bot_id":`),
			}},
		},
	})
	if err == nil {
		t.Fatal("expected malformed Twitch bot config error")
	}
}
