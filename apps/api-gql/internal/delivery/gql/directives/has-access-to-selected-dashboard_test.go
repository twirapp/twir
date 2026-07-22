package directives

import (
	"context"
	"testing"

	"github.com/google/uuid"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
)

func TestHasAccessToSelectedDashboardUsesNormalizedBindingOwnership(t *testing.T) {
	t.Parallel()

	dashboardID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	legacyOwnerID := ownerID.String()

	tests := []struct {
		name           string
		legacyChannel  model.Channels
		normalized     channelsmodel.Channel
		wantAccess     bool
		wantRoleLookup bool
	}{
		{
			name: "generic VK owner",
			normalized: channelsmodel.Channel{ID: dashboardID, Bindings: []channelplatformsmodel.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: ownerID, PlatformChannelID: "vk-channel", Enabled: true,
			}}},
			wantAccess: true,
		},
		{
			name:          "legacy owner without normalized bindings",
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			normalized:    channelsmodel.Channel{ID: dashboardID},
			wantAccess:    true,
		},
		{
			name:          "stale legacy owner denied when normalized binding belongs to another user",
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			normalized: channelsmodel.Channel{ID: dashboardID, Bindings: []channelplatformsmodel.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: otherUserID, PlatformChannelID: "vk-channel", Enabled: true,
			}}},
			wantRoleLookup: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &selectedDashboardDirectiveStore{channel: tt.legacyChannel}
			directive := &Directives{
				sessions: &selectedDashboardDirectiveSession{
					user:        &model.Users{ID: ownerID.String()},
					dashboardID: dashboardID.String(),
				},
				channels:               selectedDashboardDirectiveChannelReader{channel: tt.normalized},
				selectedDashboardStore: store,
			}
			nextCalls := 0

			result, err := directive.HasAccessToSelectedDashboard(context.Background(), nil, func(context.Context) (any, error) {
				nextCalls++
				return "allowed", nil
			})
			if tt.wantAccess {
				if err != nil {
					t.Fatalf("HasAccessToSelectedDashboard() error = %v", err)
				}
				if result != "allowed" || nextCalls != 1 {
					t.Fatalf("HasAccessToSelectedDashboard() = (%#v, next calls %d), want allowed once", result, nextCalls)
				}
			} else {
				if err == nil || err.Error() != "user does not have access to selected dashboard" {
					t.Fatalf("HasAccessToSelectedDashboard() error = %v, want access denial", err)
				}
				if nextCalls != 0 {
					t.Fatalf("HasAccessToSelectedDashboard() invoked next %d times", nextCalls)
				}
			}
			if store.roleLookups != boolToInt(tt.wantRoleLookup) {
				t.Fatalf("role lookups = %d, want %d", store.roleLookups, boolToInt(tt.wantRoleLookup))
			}
		})
	}
}

type selectedDashboardDirectiveSession struct {
	user        *model.Users
	dashboardID string
}

func (s *selectedDashboardDirectiveSession) GetAuthenticatedUserModel(context.Context) (*model.Users, error) {
	return s.user, nil
}

func (s *selectedDashboardDirectiveSession) GetAuthenticatedUserByApiKey(context.Context) (*model.Users, error) {
	return nil, context.Canceled
}

func (s *selectedDashboardDirectiveSession) GetSelectedDashboard(context.Context) (string, error) {
	return s.dashboardID, nil
}

type selectedDashboardDirectiveChannelReader struct {
	channel channelsmodel.Channel
}

func (r selectedDashboardDirectiveChannelReader) GetChannelByID(_ context.Context, channelID uuid.UUID) (channelsmodel.Channel, error) {
	if channelID != r.channel.ID {
		return channelsmodel.Nil, context.Canceled
	}
	return r.channel, nil
}

type selectedDashboardDirectiveStore struct {
	channel     model.Channels
	roleLookups int
}

func (s *selectedDashboardDirectiveStore) GetSelectedDashboardChannel(context.Context, string) (model.Channels, error) {
	return s.channel, nil
}

func (s *selectedDashboardDirectiveStore) GetSelectedDashboardRoles(context.Context, string, string) ([]model.ChannelRole, error) {
	s.roleLookups++
	return nil, nil
}

func (*selectedDashboardDirectiveStore) GetSelectedDashboardUserStat(context.Context, string, string) (model.UsersStats, error) {
	return model.UsersStats{}, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
