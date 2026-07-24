package directives

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
)

func TestHasChannelRolesDashboardPermissionUsesDashboardAccess(t *testing.T) {
	dashboardID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	legacyOwnerID := ownerID.String()
	requestedPermission := gqlmodel.ChannelRolePermissionEnum("VIEW_COMMANDS")

	tests := []struct {
		name            string
		legacyChannel   model.Channels
		normalized      channelentity.Channel
		roles           []model.ChannelRole
		permission      *gqlmodel.ChannelRolePermissionEnum
		wantAccess      bool
		wantRoleLookups int
	}{
		{
			name: "allows normalized VK-only owner",
			normalized: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: ownerID,
			}}},
			permission: &requestedPermission,
			wantAccess: true,
		},
		{
			name:          "denies stale legacy owner with a remaining binding",
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			normalized: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformKick, UserID: otherUserID,
			}}},
			permission:      &requestedPermission,
			wantRoleLookups: 1,
		},
		{
			name:       "allows assigned role with requested permission",
			normalized: channelentity.Channel{ID: dashboardID},
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: ownerID.String()}},
				Permissions: pq.StringArray{"VIEW_COMMANDS"},
			}},
			permission:      &requestedPermission,
			wantAccess:      true,
			wantRoleLookups: 1,
		},
		{
			name:       "allows assigned role with dashboard permission",
			normalized: channelentity.Channel{ID: dashboardID},
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: ownerID.String()}},
				Permissions: pq.StringArray{"CAN_ACCESS_DASHBOARD"},
			}},
			permission:      &requestedPermission,
			wantAccess:      true,
			wantRoleLookups: 1,
		},
		{
			name:       "denies nonmember role permission",
			normalized: channelentity.Channel{ID: dashboardID},
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: otherUserID.String()}},
				Permissions: pq.StringArray{"VIEW_COMMANDS"},
			}},
			permission:      &requestedPermission,
			wantRoleLookups: 1,
		},
		{
			name:          "allows zero-binding legacy owner",
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			normalized:    channelentity.Channel{ID: dashboardID},
			permission:    &requestedPermission,
			wantAccess:    true,
		},
		{
			name:       "allows assigned role with nil permission",
			normalized: channelentity.Channel{ID: dashboardID},
			roles: []model.ChannelRole{{
				Users:       []*model.ChannelRoleUser{{UserID: ownerID.String()}},
				Permissions: pq.StringArray{"VIEW_COMMANDS"},
			}},
			wantAccess:      true,
			wantRoleLookups: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &selectedDashboardDirectiveStore{
				channel: tt.legacyChannel,
				roles:   tt.roles,
			}
			directive := &Directives{
				sessions: &selectedDashboardDirectiveSession{
					user:        &model.Users{ID: ownerID.String()},
					dashboardID: dashboardID.String(),
				},
				dashboardAccess: dashboardaccess.New(
					selectedDashboardDirectiveChannelReader{channel: tt.normalized},
					store,
				),
			}
			nextCalls := 0

			result, err := directive.HasChannelRolesDashboardPermission(context.Background(), nil, func(context.Context) (any, error) {
				nextCalls++
				return "allowed", nil
			}, tt.permission)
			if tt.wantAccess {
				if err != nil {
					t.Fatalf("HasChannelRolesDashboardPermission() error = %v", err)
				}
				if result != "allowed" || nextCalls != 1 {
					t.Fatalf("HasChannelRolesDashboardPermission() = (%#v, next calls %d), want allowed once", result, nextCalls)
				}
			} else {
				if err == nil || err.Error() != "user has no permission to access this resource" {
					t.Fatalf("HasChannelRolesDashboardPermission() error = %v, want access denial", err)
				}
				if nextCalls != 0 {
					t.Fatalf("HasChannelRolesDashboardPermission() invoked next %d times", nextCalls)
				}
			}
			if store.roleLookups != tt.wantRoleLookups {
				t.Fatalf("role lookups = %d, want %d", store.roleLookups, tt.wantRoleLookups)
			}
		})
	}
}
