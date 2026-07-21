package channels

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	config "github.com/twirapp/twir/libs/config"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

func TestResolveApiKeyChannelIdentityUsesBindingLookup(t *testing.T) {
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	channelID := uuid.New()
	twitchPlatformID := "12345"
	kickPlatformID := "67890"

	tests := []struct {
		name         string
		user         usersmodel.User
		wantPlatform platformentity.Platform
		wantUserID   uuid.UUID
	}{
		{
			name:         "twitch api key resolves linked channel targets",
			user:         usersmodel.User{ID: twitchUserID, Platform: platformentity.PlatformTwitch},
			wantPlatform: platformentity.PlatformTwitch,
			wantUserID:   twitchUserID,
		},
		{
			name:         "kick api key resolves same linked channel targets",
			user:         usersmodel.User{ID: kickUserID, Platform: platformentity.PlatformKick},
			wantPlatform: platformentity.PlatformKick,
			wantUserID:   kickUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &fakeUsersRepository{user: tt.user}
			channels := &fakeChannelsRepository{channel: channelWithBindings(
				channelID,
				twitchUserID,
				twitchPlatformID,
				kickUserID,
				kickPlatformID,
			)}
			service := newTestService(users, channels)

			identity, err := service.ResolveApiKeyChannelIdentityByUserOrChannelApiKey(
				context.Background(),
				"api-key",
			)
			if err != nil {
				t.Fatalf("ResolveApiKeyChannelIdentity() error = %v", err)
			}

			if identity.InternalChannelID != channelID.String() {
				t.Fatalf("InternalChannelID = %q, want %q", identity.InternalChannelID, channelID)
			}

			gotTargets := make(map[string]string, len(identity.ChatTargets))
			for _, target := range identity.ChatTargets {
				gotTargets[target.Platform] = target.PlatformChannelID
			}

			for platform, wantID := range map[string]string{
				"twitch": twitchPlatformID,
				"kick":   kickPlatformID,
			} {
				if gotTargets[platform] != wantID {
					t.Fatalf("target %q = %q, want %q", platform, gotTargets[platform], wantID)
				}
			}

			if len(channels.bindingUserLookups) != 1 {
				t.Fatalf("GetByBindingUserID calls = %d, want 1", len(channels.bindingUserLookups))
			}
			lookup := channels.bindingUserLookups[0]
			if lookup.platform != tt.wantPlatform || lookup.userID != tt.wantUserID {
				t.Fatalf(
					"GetByBindingUserID(%s, %s), want (%s, %s)",
					lookup.platform,
					lookup.userID,
					tt.wantPlatform,
					tt.wantUserID,
				)
			}
		})
	}
}

func TestGetByPlatformIDMapsSelectedBinding(t *testing.T) {
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	channelID := uuid.New()
	channel := channelWithBindings(channelID, twitchUserID, "twitch-channel", kickUserID, "kick-channel")

	tests := []struct {
		name              string
		platform          platformentity.Platform
		platformChannelID string
		want              entity.Channel
	}{
		{
			name:              "twitch binding is selected even when it follows kick",
			platform:          platformentity.PlatformTwitch,
			platformChannelID: "twitch-channel",
			want: entity.Channel{
				ID:             channelID,
				IsEnabled:      true,
				IsTwitchBanned: true,
				IsBotMod:       true,
				BotID:          "twitch-bot",
			},
		},
		{
			name:              "kick binding state does not use twitch binding",
			platform:          platformentity.PlatformKick,
			platformChannelID: "kick-channel",
			want: entity.Channel{
				ID:        channelID,
				IsEnabled: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := &fakeChannelsRepository{channel: channel}
			service := newTestService(&fakeUsersRepository{}, channels)

			var (
				got entity.Channel
				err error
			)
			if tt.platform == platformentity.PlatformTwitch {
				got, err = service.GetByTwitchPlatformID(context.Background(), tt.platformChannelID)
			} else {
				got, err = service.GetByKickPlatformID(context.Background(), tt.platformChannelID)
			}
			if err != nil {
				t.Fatalf("GetByPlatformID() error = %v", err)
			}

			if got != tt.want {
				t.Fatalf("GetByPlatformID() = %#v, want %#v", got, tt.want)
			}
			if len(channels.platformChannelLookups) != 1 {
				t.Fatalf("GetByPlatformChannelID calls = %d, want 1", len(channels.platformChannelLookups))
			}
			lookup := channels.platformChannelLookups[0]
			if lookup.platform != tt.platform || lookup.platformChannelID != tt.platformChannelID {
				t.Fatalf(
					"GetByPlatformChannelID(%s, %q), want (%s, %q)",
					lookup.platform,
					lookup.platformChannelID,
					tt.platform,
					tt.platformChannelID,
				)
			}
		})
	}
}

func TestGetByIDUsesTwitchBindingRegardlessOfBindingOrder(t *testing.T) {
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	channelID := uuid.New()
	channels := &fakeChannelsRepository{channel: channelWithBindings(
		channelID,
		twitchUserID,
		"twitch-channel",
		kickUserID,
		"kick-channel",
	)}
	service := newTestService(&fakeUsersRepository{}, channels)

	got, err := service.GetByID(context.Background(), channelID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	want := entity.Channel{
		ID:             channelID,
		IsEnabled:      true,
		IsTwitchBanned: true,
		IsBotMod:       true,
		BotID:          "twitch-bot",
	}
	if got != want {
		t.Fatalf("GetByID() = %#v, want %#v", got, want)
	}
}

func channelWithBindings(
	channelID uuid.UUID,
	twitchUserID uuid.UUID,
	twitchPlatformID string,
	kickUserID uuid.UUID,
	kickPlatformID string,
) channelsmodel.Channel {
	return channelsmodel.Channel{
		ID: channelID,
		Bindings: []channelplatformsmodel.ChannelPlatform{
			{
				Platform:          platformentity.PlatformKick,
				UserID:            kickUserID,
				PlatformChannelID: kickPlatformID,
				Enabled:           false,
				BotConfig:         json.RawMessage(`{"kick_bot_id":"kick-bot"}`),
			},
			{
				Platform:          platformentity.PlatformTwitch,
				UserID:            twitchUserID,
				PlatformChannelID: twitchPlatformID,
				Enabled:           true,
				BotConfig: json.RawMessage(
					`{"bot_id":"twitch-bot","is_bot_mod":true,"is_twitch_banned":true}`,
				),
			},
		},
	}
}

func newTestService(users usersrepo.Repository, channels channelsrepo.Repository) *Service {
	return New(
		Opts{
			UsersRepository:    users,
			ChannelsRepository: channels,
			ChannelService: channelservice.NewChannelService(
				channels,
				&buscore.Bus{},
				config.Config{},
				nil,
				nil,
			),
		},
	)
}

type fakeUsersRepository struct {
	user usersmodel.User
}

func (f *fakeUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return f.user, nil
}

func (f *fakeUsersRepository) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	return f.user, nil
}

func (f *fakeUsersRepository) GetByPlatformID(context.Context, platformentity.Platform, string) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (f *fakeUsersRepository) GetManyByIDS(context.Context, usersrepo.GetManyInput) ([]usersmodel.User, error) {
	return nil, nil
}

func (f *fakeUsersRepository) Update(context.Context, uuid.UUID, usersrepo.UpdateInput) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

func (f *fakeUsersRepository) GetRandomOnlineUser(context.Context, usersrepo.GetRandomOnlineUserInput) (usersmodel.OnlineUser, error) {
	return usersmodel.NilOnlineUser, nil
}

func (f *fakeUsersRepository) GetOnlineUsersWithFilters(context.Context, usersrepo.GetOnlineUsersWithFiltersInput) ([]usersmodel.OnlineUser, error) {
	return nil, nil
}

func (f *fakeUsersRepository) Create(context.Context, usersrepo.CreateInput) (usersmodel.User, error) {
	return usersmodel.Nil, nil
}

type fakeChannelsRepository struct {
	channel                channelsmodel.Channel
	bindingUserLookups     []bindingUserLookup
	platformChannelLookups []platformChannelLookup
}

type bindingUserLookup struct {
	platform platformentity.Platform
	userID   uuid.UUID
}

type platformChannelLookup struct {
	platform          platformentity.Platform
	platformChannelID string
}

func (f *fakeChannelsRepository) GetByBindingUserID(
	_ context.Context,
	p platformentity.Platform,
	userID uuid.UUID,
) (channelsmodel.Channel, error) {
	f.bindingUserLookups = append(f.bindingUserLookups, bindingUserLookup{platform: p, userID: userID})
	return f.channel, nil
}

func (f *fakeChannelsRepository) GetByPlatformChannelID(
	_ context.Context,
	p platformentity.Platform,
	platformChannelID string,
) (channelsmodel.Channel, error) {
	f.platformChannelLookups = append(
		f.platformChannelLookups,
		platformChannelLookup{platform: p, platformChannelID: platformChannelID},
	)
	return f.channel, nil
}

func (f *fakeChannelsRepository) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (f *fakeChannelsRepository) GetByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return f.channel, nil
}

func (f *fakeChannelsRepository) GetBySlug(context.Context, channelsrepo.GetBySlugInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) GetCount(context.Context, channelsrepo.GetCountInput) (int, error) {
	return 0, nil
}

func (f *fakeChannelsRepository) Update(context.Context, uuid.UUID, channelsrepo.UpdateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) Create(context.Context, channelsrepo.CreateInput) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}

func (f *fakeChannelsRepository) GetByApiKey(context.Context, string) (channelsmodel.Channel, error) {
	return channelsmodel.Nil, nil
}
