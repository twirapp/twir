package resolvers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelsrepo "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func TestResolveApiKeyChannelIdentity(t *testing.T) {
	twitchUserID := uuid.New()
	kickUserID := uuid.New()
	channelID := uuid.New()
	twitchPlatformID := "12345"
	kickPlatformID := "67890"

	tests := []struct {
		name            string
		user            usersmodel.User
		wantChannelID   string
		wantTargets     map[string]string
		wantTwitchCalls int
		wantKickCalls   int
	}{
		{
			name:          "twitch api key resolves linked channel targets",
			user:          usersmodel.User{ID: twitchUserID, Platform: platformentity.PlatformTwitch},
			wantChannelID: channelID.String(),
			wantTargets: map[string]string{
				"twitch": twitchPlatformID,
				"kick":   kickPlatformID,
			},
			wantTwitchCalls: 1,
		},
		{
			name:          "kick api key resolves same linked channel targets",
			user:          usersmodel.User{ID: kickUserID, Platform: platformentity.PlatformKick},
			wantChannelID: channelID.String(),
			wantTargets: map[string]string{
				"twitch": twitchPlatformID,
				"kick":   kickPlatformID,
			},
			wantKickCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &fakeUsersRepository{user: tt.user}
			channels := &fakeChannelsRepository{channel: channelsmodel.Channel{
				ID:               channelID,
				TwitchUserID:     &twitchUserID,
				TwitchPlatformID: &twitchPlatformID,
				KickUserID:       &kickUserID,
				KickPlatformID:   &kickPlatformID,
			}}

			identity, err := resolveApiKeyChannelIdentity(context.Background(), Deps{
				UsersRepository:    users,
				ChannelsRepository: channels,
			}, "api-key")
			if err != nil {
				t.Fatalf("resolveApiKeyChannelIdentity() error = %v", err)
			}

			if identity.InternalChannelID != tt.wantChannelID {
				t.Fatalf("InternalChannelID = %q, want %q", identity.InternalChannelID, tt.wantChannelID)
			}

			gotTargets := make(map[string]string, len(identity.ChatTargets))
			for _, target := range identity.ChatTargets {
				gotTargets[target.Platform] = target.PlatformChannelID
			}

			for platform, wantID := range tt.wantTargets {
				if gotTargets[platform] != wantID {
					t.Fatalf("target %q = %q, want %q", platform, gotTargets[platform], wantID)
				}
			}

			if channels.twitchCalls != tt.wantTwitchCalls {
				t.Fatalf("GetByTwitchUserID calls = %d, want %d", channels.twitchCalls, tt.wantTwitchCalls)
			}
			if channels.kickCalls != tt.wantKickCalls {
				t.Fatalf("GetByKickUserID calls = %d, want %d", channels.kickCalls, tt.wantKickCalls)
			}
		})
	}
}

type fakeUsersRepository struct {
	user usersmodel.User
}

func (f *fakeUsersRepository) GetByApiKey(context.Context, string) (usersmodel.User, error) {
	return f.user, nil
}

func (f *fakeUsersRepository) GetByID(context.Context, uuid.UUID) (usersmodel.User, error) {
	return usersmodel.Nil, nil
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
	channel     channelsmodel.Channel
	twitchCalls int
	kickCalls   int
}

func (f *fakeChannelsRepository) GetByTwitchUserID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	f.twitchCalls++
	return f.channel, nil
}

func (f *fakeChannelsRepository) GetByKickUserID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	f.kickCalls++
	return f.channel, nil
}

func (f *fakeChannelsRepository) GetMany(context.Context, channelsrepo.GetManyInput) ([]channelsmodel.Channel, error) {
	return nil, nil
}

func (f *fakeChannelsRepository) GetByID(context.Context, uuid.UUID) (channelsmodel.Channel, error) {
	return f.channel, nil
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
