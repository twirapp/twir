package twir_stats

import (
	"context"
	"reflect"
	"strings"
	"testing"

	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestVisibleChannelBindingsQueryExcludesBannedTwitchBindings(t *testing.T) {
	t.Parallel()

	twitch := platformentity.PlatformTwitch
	tests := []struct {
		name     string
		platform *platformentity.Platform
		wantVars []any
	}{
		{
			name:     "all channels",
			wantVars: []any{true, false, platformentity.PlatformTwitch, `{"is_twitch_banned":true}`},
		},
		{
			name:     "Twitch channels",
			platform: &twitch,
			wantVars: []any{true, false, platformentity.PlatformTwitch, `{"is_twitch_banned":true}`, platformentity.PlatformTwitch},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statement := visibleChannelBindingsQuery(
				newStatsDryRunPostgresDB(t),
				context.Background(),
				tt.platform,
			).Distinct("cp.channel_id").Count(new(int64)).Statement
			query := statement.SQL.String()

			for _, fragment := range []string{
				"FROM channel_platforms AS cp",
				"JOIN users AS u ON u.id = cp.user_id",
				"cp.enabled = $1",
				"u.is_banned = $2",
				"NOT (cp.platform = $3 AND cp.bot_config @> $4::jsonb)",
			} {
				if !strings.Contains(query, fragment) {
					t.Fatalf("query does not contain %q: %s", fragment, query)
				}
			}
			if tt.platform != nil && !strings.Contains(query, "cp.platform = $5") {
				t.Fatalf("Twitch query does not filter its platform: %s", query)
			}
			if !reflect.DeepEqual(statement.Vars, tt.wantVars) {
				t.Fatalf("query vars = %#v, want %#v", statement.Vars, tt.wantVars)
			}
		})
	}
}

func newStatsDryRunPostgresDB(t *testing.T) *gorm.DB {
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
