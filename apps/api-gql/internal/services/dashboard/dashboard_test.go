package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	apiChannelbinding "github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

type dashboardCurrentPlatformStub struct {
	platform string
	err      error
}

func (s dashboardCurrentPlatformStub) GetCurrentPlatform(context.Context) (string, error) {
	return s.platform, s.err
}

type dashboardChannelLookupStub struct {
	channel channelsmodel.Channel
}

func (s dashboardChannelLookupStub) GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return s.channel, nil
}

type dashboardUsersLookupStub struct {
	users map[uuid.UUID]usersmodel.User
}

func (s dashboardUsersLookupStub) GetByID(_ context.Context, id uuid.UUID) (usersmodel.User, error) {
	user, ok := s.users[id]
	if !ok {
		return usersmodel.Nil, errors.New("user not found")
	}

	return user, nil
}

func TestResolveAnalyticsIdentitySelectsCurrentPlatformBinding(t *testing.T) {
	service := &Service{authService: dashboardCurrentPlatformStub{platform: platform.PlatformKick.String()}}
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformTwitch, PlatformChannelID: "twitch-channel"},
			{Platform: platform.PlatformKick, PlatformChannelID: "kick-channel"},
		},
	}

	gotPlatform, gotChannelID := service.resolveAnalyticsIdentity(context.Background(), channel)
	if gotPlatform != platform.PlatformKick.String() || gotChannelID != "kick-channel" {
		t.Fatalf("analytics identity = (%q, %q), want Kick binding", gotPlatform, gotChannelID)
	}
}

func TestResolveAnalyticsIdentityFallsBackToTwitchBinding(t *testing.T) {
	service := &Service{authService: dashboardCurrentPlatformStub{err: errors.New("no current platform")}}
	channel := channelsmodel.Channel{
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformKick, PlatformChannelID: "kick-channel"},
			{Platform: platform.PlatformTwitch, PlatformChannelID: "twitch-channel"},
		},
	}

	gotPlatform, gotChannelID := service.resolveAnalyticsIdentity(context.Background(), channel)
	if gotPlatform != platform.PlatformTwitch.String() || gotChannelID != "twitch-channel" {
		t.Fatalf("analytics identity = (%q, %q), want Twitch fallback", gotPlatform, gotChannelID)
	}
}

func TestGetBotStatusesMapsKickBindingIdentity(t *testing.T) {
	channelID := uuid.New()
	kickOwnerID := uuid.New()
	kickBotUserID := uuid.New()
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelsmodel.Channel{
			ID: channelID,
			Bindings: []channelplatformsmodel.ChannelPlatform{{
				Platform:          platform.PlatformKick,
				UserID:            kickOwnerID,
				PlatformChannelID: "kick-channel",
				Enabled:           true,
				BotUserID:         &kickBotUserID,
			}},
		}},
		usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
			kickOwnerID:   {ID: kickOwnerID, Login: "kick-owner"},
			kickBotUserID: {ID: kickBotUserID, Login: "kick-bot"},
		}},
	}

	statuses, err := service.GetBotStatuses(context.Background(), channelID.String())
	if err != nil {
		t.Fatalf("get bot statuses: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	status := statuses[0]
	if status.Platform != platform.PlatformKick.String() || status.ChannelName != "kick-owner" {
		t.Fatalf("Kick status = %#v, want binding owner identity", status)
	}
	if !status.Enabled || !status.IsMod {
		t.Fatalf("Kick status = %#v, want enabled moderator", status)
	}
	if status.BotID != kickBotUserID.String() || status.BotName != "kick-bot" {
		t.Fatalf("Kick status = %#v, want binding bot identity", status)
	}
}

func TestGetBasicTwitchBotStatusUsesTwitchBindingConfig(t *testing.T) {
	twitchOwnerID := uuid.New()
	channel := channelsmodel.Channel{
		ID: uuid.New(),
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{Platform: platform.PlatformKick, Enabled: false},
			{
				Platform:          platform.PlatformTwitch,
				UserID:            twitchOwnerID,
				PlatformChannelID: "twitch-channel",
				Enabled:           true,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
		},
	}
	binding, config, found, err := apiChannelbinding.FindTwitch(channel)
	if err != nil {
		t.Fatalf("find Twitch binding: %v", err)
	}
	if !found {
		t.Fatal("expected Twitch binding")
	}

	service := &Service{usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
		twitchOwnerID: {ID: twitchOwnerID, Login: "twitch-owner"},
	}}}
	status := service.getBasicTwitchBotStatus(context.Background(), channel, binding, config)
	if status.Platform != platform.PlatformTwitch.String() || status.ChannelName != "twitch-owner" {
		t.Fatalf("Twitch status = %#v, want binding owner identity", status)
	}
	if !status.Enabled || !status.IsMod || status.BotID != "twitch-bot" || status.BotName != "TwirBot" {
		t.Fatalf("Twitch status = %#v, want parsed Twitch config", status)
	}
}
