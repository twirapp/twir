package resolvers

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/gorm"
)

func TestCommunityResetStatsUsesNormalizedOwnerCheck(t *testing.T) {
	dashboardID := uuid.New()
	ownerID := uuid.New()
	otherUserID := uuid.New()
	legacyOwnerID := ownerID.String()

	tests := []struct {
		name          string
		normalized    channelentity.Channel
		legacyChannel model.Channels
		wantSuccess   bool
		wantUpdates   int
	}{
		{
			name: "denies stale legacy owner without an update",
			normalized: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: otherUserID,
			}}},
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
		},
		{
			name: "allows normalized owner",
			normalized: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: ownerID,
			}}},
			wantSuccess: true,
			wantUpdates: 1,
		},
		{
			name:          "allows zero-binding legacy owner",
			normalized:    channelentity.Channel{ID: dashboardID},
			legacyChannel: model.Channels{TwitchUserID: &legacyOwnerID},
			wantSuccess:   true,
			wantUpdates:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newResolverDryRunPostgresDB(t).Session(&gorm.Session{SkipDefaultTransaction: true})
			updates := 0
			if err := db.Callback().Update().Before("gorm:update").Register("community-reset-stats-count-update", func(*gorm.DB) {
				updates++
			}); err != nil {
				t.Fatalf("register update callback: %v", err)
			}

			resolver := &mutationResolver{Resolver: &Resolver{deps: Deps{
				Gorm: db,
				DashboardAccess: dashboardaccess.New(
					communityResetStatsChannelReader{channel: tt.normalized},
					&communityResetStatsStore{legacyChannel: tt.legacyChannel},
				),
			}}}

			got, err := resolver.resetCommunityStats(
				context.Background(),
				ownerID.String(),
				dashboardID.String(),
				gqlmodel.CommunityUsersResetTypeMessages,
			)
			if tt.wantSuccess {
				if err != nil {
					t.Fatalf("resetCommunityStats() error = %v", err)
				}
				if !got {
					t.Fatal("resetCommunityStats() = false, want true")
				}
			} else {
				if err == nil {
					t.Fatal("resetCommunityStats() error = nil, want owner denial")
				}
				if got {
					t.Fatal("resetCommunityStats() = true, want false")
				}
			}
			if updates != tt.wantUpdates {
				t.Fatalf("stats updates = %d, want %d", updates, tt.wantUpdates)
			}
		})
	}
}

type communityResetStatsChannelReader struct {
	channel channelentity.Channel
}

func (r communityResetStatsChannelReader) GetChannelByID(_ context.Context, channelID uuid.UUID) (channelentity.Channel, error) {
	if channelID != r.channel.ID {
		return channelentity.Nil, context.Canceled
	}

	return r.channel, nil
}

type communityResetStatsStore struct {
	legacyChannel model.Channels
}

func (s *communityResetStatsStore) GetLegacyChannel(context.Context, uuid.UUID) (model.Channels, error) {
	return s.legacyChannel, nil
}

func (*communityResetStatsStore) GetRoles(context.Context, uuid.UUID, string) ([]model.ChannelRole, error) {
	return nil, nil
}

func (*communityResetStatsStore) GetUserStat(context.Context, string, uuid.UUID) (model.UsersStats, error) {
	return model.UsersStats{}, nil
}
