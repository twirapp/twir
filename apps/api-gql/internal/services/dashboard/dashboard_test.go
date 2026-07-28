package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
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
	channel channelentity.Channel
}

func (s dashboardChannelLookupStub) GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error) {
	return s.channel, nil
}

type dashboardUsersLookupStub struct {
	users map[uuid.UUID]usersmodel.User
}

type dashboardBindingPatch struct {
	id    uuid.UUID
	input channelplatforms.PatchInput
}

type dashboardBindingUpdaterStub struct {
	patches []dashboardBindingPatch
}

func (s *dashboardBindingUpdaterStub) Patch(
	_ context.Context,
	id uuid.UUID,
	input channelplatforms.PatchInput,
) (channelplatformentity.ChannelPlatform, error) {
	s.patches = append(s.patches, dashboardBindingPatch{id: id, input: input})
	return channelplatformentity.ChannelPlatform{ID: id}, nil
}

type dashboardCacheInvalidatorStub struct {
	keys []string
}

func (s *dashboardCacheInvalidatorStub) Invalidate(_ context.Context, key string) error {
	s.keys = append(s.keys, key)
	return nil
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
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{
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
	channel := channelentity.Channel{
		Bindings: []channelplatformentity.ChannelPlatform{
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
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
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

func TestGetBotStatusesMapsVKVideoLiveBindingIdentity(t *testing.T) {
	// Given
	channelID := uuid.New()
	vkOwnerID := uuid.New()
	vkBotUserID := uuid.New()
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				Platform:          platform.PlatformVKVideoLive,
				UserID:            vkOwnerID,
				PlatformChannelID: "vk-channel",
				Enabled:           true,
				BotUserID:         &vkBotUserID,
			}},
		}},
		usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
			vkOwnerID:   {ID: vkOwnerID, Login: "vk-owner"},
			vkBotUserID: {ID: vkBotUserID, Login: "vk-bot"},
		}},
	}

	// When
	statuses, err := service.GetBotStatuses(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get bot statuses: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	status := statuses[0]
	if status.Platform != platform.PlatformVKVideoLive.String() || status.ChannelName != "vk-owner" {
		t.Fatalf("VK Video Live status = %#v, want binding owner identity", status)
	}
	if !status.Enabled {
		t.Fatalf("VK Video Live status = %#v, want enabled binding", status)
	}
	if status.BotID != vkBotUserID.String() || status.BotName != "vk-bot" {
		t.Fatalf("VK Video Live status = %#v, want global bot identity", status)
	}
}

func TestGetBotStatusesDoesNotExposeUsableVKVideoLiveBotWithoutBotUserID(t *testing.T) {
	// Given
	channelID := uuid.New()
	vkOwnerID := uuid.New()
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				Platform:          platform.PlatformVKVideoLive,
				UserID:            vkOwnerID,
				PlatformChannelID: "vk-channel",
				Enabled:           true,
			}},
		}},
		usersRepo: dashboardUsersLookupStub{users: map[uuid.UUID]usersmodel.User{
			vkOwnerID: {ID: vkOwnerID, Login: "vk-owner"},
		}},
	}

	// When
	statuses, err := service.GetBotStatuses(context.Background(), channelID.String())

	// Then
	if err != nil {
		t.Fatalf("get bot statuses: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses = %d, want 1", len(statuses))
	}
	status := statuses[0]
	if status.Platform != platform.PlatformVKVideoLive.String() || status.ChannelName != "vk-owner" {
		t.Fatalf("VK Video Live status = %#v, want binding owner identity", status)
	}
	if status.BotID != "" || status.BotName != "" {
		t.Fatalf("VK Video Live status = %#v, want no usable bot identity", status)
	}
}

func TestGetBasicTwitchBotStatusUsesTwitchBindingConfig(t *testing.T) {
	twitchOwnerID := uuid.New()
	channel := channelentity.Channel{
		ID: uuid.New(),
		Bindings: []channelplatformentity.ChannelPlatform{
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
	binding, config, found, err := channel.TwitchBinding()
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

func TestBotJoinLeaveUpdatesOnlyVKVideoLiveBindingWithoutEventSub(t *testing.T) {
	for _, test := range []struct {
		name        string
		action      string
		platform    string
		includeKick bool
		wantEnabled bool
	}{
		{name: "joins", action: BotJoinLeaveActionJoin, platform: platform.PlatformVKVideoLive.String(), includeKick: true, wantEnabled: true},
		{name: "leaves", action: BotJoinLeaveActionLeave, platform: platform.PlatformVKVideoLive.String(), includeKick: true, wantEnabled: false},
		{name: "defaults to VK binding", action: BotJoinLeaveActionJoin, wantEnabled: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			// Given
			channelID := uuid.New()
			vkBindingID := uuid.New()
			vkBotUserID := uuid.New()
			updater := &dashboardBindingUpdaterStub{}
			cache := &dashboardCacheInvalidatorStub{}
			bindings := []channelplatformentity.ChannelPlatform{
				{ID: vkBindingID, Platform: platform.PlatformVKVideoLive, PlatformChannelID: "vk-channel", BotUserID: &vkBotUserID},
			}
			if test.includeKick {
				bindings = append(bindings, channelplatformentity.ChannelPlatform{ID: uuid.New(), Platform: platform.PlatformKick, Enabled: true})
			}
			service := &Service{
				channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
					ID:       channelID,
					Bindings: bindings,
				}},
				channelPlatformsRepo: updater,
				channelsCache:        cache,
			}

			// When
			success, err := service.BotJoinLeave(context.Background(), channelID.String(), test.action, test.platform)

			// Then
			if err != nil {
				t.Fatalf("bot join/leave: %v", err)
			}
			if !success {
				t.Fatal("bot join/leave = false, want true")
			}
			if len(updater.patches) != 1 {
				t.Fatalf("patches = %d, want 1", len(updater.patches))
			}
			patch := updater.patches[0]
			if patch.id != vkBindingID || patch.input.Enabled == nil || *patch.input.Enabled != test.wantEnabled {
				t.Fatalf("patch = %#v, want VK binding enabled=%t", patch, test.wantEnabled)
			}
			if len(cache.keys) != 1 || cache.keys[0] != channelID.String() {
				t.Fatalf("cache invalidations = %#v, want [%q]", cache.keys, channelID.String())
			}
		})
	}
}

func TestBotJoinLeaveRejectsVKVideoLiveBindingWithoutBotUserID(t *testing.T) {
	// Given
	channelID := uuid.New()
	updater := &dashboardBindingUpdaterStub{}
	service := &Service{
		channelService: dashboardChannelLookupStub{channel: channelentity.Channel{
			ID: channelID,
			Bindings: []channelplatformentity.ChannelPlatform{{
				ID:                uuid.New(),
				Platform:          platform.PlatformVKVideoLive,
				PlatformChannelID: "vk-channel",
			}},
		}},
		channelPlatformsRepo: updater,
	}

	// When
	success, err := service.BotJoinLeave(context.Background(), channelID.String(), BotJoinLeaveActionJoin, platform.PlatformVKVideoLive.String())

	// Then
	if err == nil {
		t.Fatal("bot join/leave error = nil, want missing VK bot error")
	}
	if success {
		t.Fatal("bot join/leave = true, want false")
	}
	if len(updater.patches) != 0 {
		t.Fatalf("patches = %d, want 0", len(updater.patches))
	}
}
