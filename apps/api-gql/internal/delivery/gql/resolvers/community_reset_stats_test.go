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

	tests := []struct {
		name        string
		normalized  channelentity.Channel
		wantSuccess bool
		wantUpdates int
	}{
		{
			name: "denies non-owner when binding belongs to another user",
			normalized: channelentity.Channel{ID: dashboardID, Bindings: []channelplatformentity.ChannelPlatform{{
				ID: uuid.New(), ChannelID: dashboardID, Platform: platformentity.PlatformVKVideoLive, UserID: otherUserID,
			}}},
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
			name:       "denies owner without bindings",
			normalized: channelentity.Channel{ID: dashboardID},
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
					&communityResetStatsStore{},
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

type communityResetStatsStore struct{}

func (*communityResetStatsStore) GetRoles(context.Context, uuid.UUID, string) ([]model.ChannelRole, error) {
	return nil, nil
}

func (*communityResetStatsStore) GetUserStat(context.Context, string, uuid.UUID) (model.UsersStats, error) {
	return model.UsersStats{}, nil
}
