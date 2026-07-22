package middlewares

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/enums/dashboard_permissions"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

type fakeChannelGetter struct {
	channel channelsmodel.Channel
	gotID   uuid.UUID
}

func (f *fakeChannelGetter) GetChannelByID(_ context.Context, channelID uuid.UUID) (channelsmodel.Channel, error) {
	f.gotID = channelID
	return f.channel, nil
}

func TestMiddlewaresIsSelectedDashboardOwnerUsesNormalizedBindings(t *testing.T) {
	dashboardID := uuid.New()
	twitchOwnerID := uuid.New()
	kickOwnerID := uuid.New()
	vkVideoLiveOwnerID := uuid.New()

	bindings := []channelplatformsmodel.ChannelPlatform{
		{Platform: platform.PlatformVKVideoLive, UserID: vkVideoLiveOwnerID},
		{Platform: platform.PlatformKick, UserID: kickOwnerID},
		{Platform: platform.PlatformTwitch, UserID: twitchOwnerID},
	}

	tests := []struct {
		name   string
		userID string
		want   bool
	}{
		{name: "Twitch owner", userID: twitchOwnerID.String(), want: true},
		{name: "Kick owner", userID: kickOwnerID.String(), want: true},
		{name: "VK Video Live owner", userID: vkVideoLiveOwnerID.String(), want: true},
		{name: "unlinked user", userID: uuid.NewString(), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getter := &fakeChannelGetter{
				channel: channelsmodel.Channel{ID: dashboardID, Bindings: bindings},
			}
			middlewares := &Middlewares{channelGetter: getter}

			got, err := middlewares.isSelectedDashboardOwner(context.Background(), dashboardID.String(), tt.userID)
			if err != nil {
				t.Fatalf("isSelectedDashboardOwner returned error: %v", err)
			}
			if got != tt.want {
				t.Errorf("isSelectedDashboardOwner() = %v, want %v", got, tt.want)
			}
			if getter.gotID != dashboardID {
				t.Errorf("GetChannelByID received %s, want %s", getter.gotID, dashboardID)
			}
		})
	}
}

func TestHasChannelRolesDashboardAccess(t *testing.T) {
	userID := uuid.NewString()
	manageCommands := dashboard_permissions.ChannelRolePermissionEnumManageCommands

	tests := []struct {
		name       string
		roles      []model.ChannelRole
		userStats  model.UsersStats
		permission *dashboard_permissions.ChannelRolePermissionEnum
		want       bool
	}{
		{
			name: "allows any dashboard access permission for an assigned role",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{dashboard_permissions.ChannelRolePermissionEnumViewCommands.String()},
			}},
			want: true,
		},
		{
			name: "allows the requested permission for an assigned role",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{manageCommands.String()},
			}},
			permission: &manageCommands,
			want:       true,
		},
		{
			name: "allows the dashboard access permission for an assigned role",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{dashboard_permissions.ChannelRolePermissionEnumCanAccessDashboard.String()},
			}},
			permission: &manageCommands,
			want:       true,
		},
		{
			name: "allows a moderator role for a moderator user",
			roles: []model.ChannelRole{{
				Type:        model.ChannelRoleTypeModerator,
				Permissions: pq.StringArray{dashboard_permissions.ChannelRolePermissionEnumViewCommands.String()},
			}},
			userStats: model.UsersStats{IsMod: true},
			want:      true,
		},
		{
			name: "denies a role without the requested permission",
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: userID}},
				Permissions: pq.StringArray{dashboard_permissions.ChannelRolePermissionEnumViewCommands.String()},
			}},
			permission: &manageCommands,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasChannelRolesDashboardAccess(tt.roles, userID, tt.userStats, tt.permission); got != tt.want {
				t.Errorf("hasChannelRolesDashboardAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}
