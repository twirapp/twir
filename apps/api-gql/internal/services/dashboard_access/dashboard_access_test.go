package dashboard_access

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
)

func TestCanAccessUsesNormalizedOwnershipAndLegacyFallback(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	legacyOwnerID := ownerID.String()

	tests := []struct {
		name          string
		channel       channelentity.Channel
		legacyChannel model.Channels
		want          bool
	}{
		{
			name: "allows a normalized binding owner",
			channel: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: ownerID,
			}}},
			want: true,
		},
		{
			name:          "allows a legacy owner when no bindings exist",
			channel:       channelentity.Channel{ID: dashboardID},
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			want:          true,
		},
		{
			name: "denies stale legacy ownership when bindings exist",
			channel: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: otherUserID,
			}}},
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(
				testChannelReader{channel: tt.channel},
				&testStore{channel: tt.legacyChannel},
			)

			got, err := service.CanAccess(
				context.Background(),
				Subject{ID: ownerID.String()},
				dashboardID,
				"",
			)
			if err != nil {
				t.Fatalf("CanAccess() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("CanAccess() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestIsOwnerUsesNormalizedOwnershipAndLegacyFallback(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	legacyOwnerID := ownerID.String()

	tests := []struct {
		name          string
		channel       channelentity.Channel
		legacyChannel model.Channels
		want          bool
	}{
		{
			name: "allows a normalized VK binding owner",
			channel: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: ownerID,
			}}},
			want: true,
		},
		{
			name:          "allows a legacy owner when no bindings exist",
			channel:       channelentity.Channel{ID: dashboardID},
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			want:          true,
		},
		{
			name: "denies stale legacy ownership when a binding remains",
			channel: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: otherUserID,
			}}},
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(
				testChannelReader{channel: tt.channel},
				&testStore{channel: tt.legacyChannel},
			)

			got, err := service.IsOwner(context.Background(), ownerID.String(), dashboardID)
			if err != nil {
				t.Fatalf("IsOwner() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("IsOwner() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestCanAccessAllowsDashboardPermissionForRole(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	userID := uuid.NewString()

	tests := []struct {
		name       string
		roles      []model.ChannelRole
		userStats  model.UsersStats
		permission string
		want       bool
	}{
		{
			name: "allows any permission for an assigned role",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{"VIEW_COMMANDS"},
			}},
			want: true,
		},
		{
			name: "allows the requested permission for an assigned role",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{"MANAGE_COMMANDS"},
			}},
			permission: "MANAGE_COMMANDS",
			want:       true,
		},
		{
			name: "allows dashboard permission for an assigned role",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{"CAN_ACCESS_DASHBOARD"},
			}},
			permission: "MANAGE_COMMANDS",
			want:       true,
		},
		{
			name: "allows a moderator role for a moderator user",
			roles: []model.ChannelRole{{
				Type:        model.ChannelRoleTypeModerator,
				Permissions: pq.StringArray{"VIEW_COMMANDS"},
			}},
			userStats: model.UsersStats{IsMod: true},
			want:      true,
		},
		{
			name: "denies a role without the requested permission",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{"VIEW_COMMANDS"},
			}},
			permission: "MANAGE_COMMANDS",
			want:       false,
		},
		{
			name: "denies permissions assigned to a different user",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: uuid.NewString()}},
				Permissions: pq.StringArray{"CAN_ACCESS_DASHBOARD", "MANAGE_COMMANDS"},
			}},
			permission: "MANAGE_COMMANDS",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(
				testChannelReader{channel: channelentity.Channel{ID: dashboardID}},
				&testStore{roles: tt.roles, stat: tt.userStats},
			)

			got, err := service.CanAccess(
				context.Background(),
				Subject{ID: userID},
				dashboardID,
				tt.permission,
			)
			if err != nil {
				t.Fatalf("CanAccess() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("CanAccess() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestCanAccessAllowsBotAdminWithoutDependencies(t *testing.T) {
	t.Parallel()

	got, err := New(nil, nil).CanAccess(
		context.Background(),
		Subject{IsBotAdmin: true},
		uuid.New(),
		"MANAGE_COMMANDS",
	)
	if err != nil {
		t.Fatalf("CanAccess() error = %v", err)
	}
	if !got {
		t.Fatal("CanAccess() = false, want bot admin access")
	}
}

type testChannelReader struct {
	channel channelentity.Channel
}

func (r testChannelReader) GetChannelByID(_ context.Context, channelID uuid.UUID) (channelentity.Channel, error) {
	if channelID != r.channel.ID {
		return channelentity.Nil, context.Canceled
	}

	return r.channel, nil
}

type testStore struct {
	channel model.Channels
	roles   []model.ChannelRole
	stat    model.UsersStats
}

func (s *testStore) GetLegacyChannel(context.Context, uuid.UUID) (model.Channels, error) {
	return s.channel, nil
}

func (s *testStore) GetRoles(context.Context, uuid.UUID, string) ([]model.ChannelRole, error) {
	return s.roles, nil
}

func (s *testStore) GetUserStat(context.Context, string, uuid.UUID) (model.UsersStats, error) {
	return s.stat, nil
}
