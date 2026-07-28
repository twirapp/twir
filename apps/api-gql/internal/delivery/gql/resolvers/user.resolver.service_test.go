package resolvers

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestOwnedDashboardsQueryUsesOwnershipBindings(t *testing.T) {
	t.Parallel()

	const userID = "a8ff9a5f-b450-44eb-ba6c-5047888e3af7"
	statement := ownedDashboardsQuery(
		newResolverDryRunPostgresDB(t),
		context.Background(),
		userID,
	).Find(&[]model.Channels{}).Statement
	query := statement.SQL.String()

	for _, fragment := range []string{
		"EXISTS (",
		"FROM channel_platforms AS cp_owner",
		"cp_owner.channel_id = channels.id",
		"cp_owner.user_id = $1::uuid",
	} {
		if !strings.Contains(query, fragment) {
			t.Fatalf("query does not contain %q: %s", fragment, query)
		}
	}

	for _, absent := range []string{"twitch_user_id", "kick_user_id", "NOT EXISTS"} {
		if strings.Contains(query, absent) {
			t.Fatalf("query must not contain %q: %s", absent, query)
		}
	}

	if want := []any{userID}; !reflect.DeepEqual(statement.Vars, want) {
		t.Fatalf("query vars = %#v, want %#v", statement.Vars, want)
	}
}

func TestResolveDashboardPlatform(t *testing.T) {
	t.Parallel()

	channelID := uuid.New()
	userID := uuid.NewString()

	tests := []struct {
		name      string
		channelID string
		reader    resolverDashboardPlatformReader
		want      string
	}{
		{
			name:      "zero binding channel has no platform",
			channelID: channelID.String(),
			reader:    resolverDashboardPlatformReader{channel: channelentity.Channel{ID: channelID}},
			want:      "",
		},
		{
			name:      "binding owner gets binding platform",
			channelID: channelID.String(),
			reader: resolverDashboardPlatformReader{channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{{
				Platform: platformentity.PlatformVKVideoLive, UserID: uuid.MustParse(userID),
			}}}},
			want: "vk_video_live",
		},
		{
			name:      "non-owner gets first binding platform",
			channelID: channelID.String(),
			reader: resolverDashboardPlatformReader{channel: channelentity.Channel{ID: channelID, Bindings: []channelplatformentity.ChannelPlatform{{
				Platform: platformentity.PlatformVKVideoLive, UserID: uuid.New(),
			}}}},
			want: "vk_video_live",
		},
		{
			name:      "invalid channel ID resolves to empty platform",
			channelID: "not-a-uuid",
			want:      "",
		},
		{
			name:      "normalized lookup error resolves to empty platform",
			channelID: channelID.String(),
			reader:    resolverDashboardPlatformReader{err: errors.New("lookup failed")},
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveDashboardPlatform(context.Background(), tt.reader, tt.channelID, userID)
			if got != tt.want {
				t.Fatalf("resolveDashboardPlatform() = %q, want %q", got, tt.want)
			}
		})
	}
}

type resolverDashboardPlatformReader struct {
	channel channelentity.Channel
	err     error
}

func (r resolverDashboardPlatformReader) GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error) {
	return r.channel, r.err
}

func newResolverDryRunPostgresDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(
		postgres.Open("host=127.0.0.1 user=twir dbname=twir sslmode=disable"),
		&gorm.Config{DisableAutomaticPing: true, DryRun: true},
	)
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}

	return db
}
