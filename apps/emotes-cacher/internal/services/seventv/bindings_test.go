package seventv

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/twirapp/twir/libs/entities/platform"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBuildStartupChannelsQueryUsesEnabledTwitchOrKickBindingWithoutDuplicateChannels(t *testing.T) {
	t.Parallel()

	db := newDryRunPostgresDB(t)
	statement := buildStartupChannelsQuery(db, context.Background()).Find(&[]string{}).Statement
	query := statement.SQL.String()

	assertQueryContains(t, query,
		"SELECT c.id FROM channels AS c",
		"EXISTS (",
		"FROM channel_platforms cp",
		"cp.channel_id = c.id",
		"cp.platform IN ($1, $2)",
		"cp.enabled = $3",
	)
	assertQueryExcludes(t, query,
		"twitch_bot_enabled",
		"kick_bot_enabled",
		"twitch_user_id",
		"kick_user_id",
		`c."isEnabled"`,
	)
	assertSingleBindingSource(t, query)
	assertQueryVars(t, statement.Vars, platform.PlatformTwitch, platform.PlatformKick, true)
}

func TestBuildChannelBindingsQueryUsesExplicitTwitchAndKickBindingsRegardlessOfBindingEnablement(t *testing.T) {
	t.Parallel()

	db := newDryRunPostgresDB(t)
	statement := buildChannelBindingsQuery(db, context.Background(), "channel-id").Find(&[]channelBindingData{}).Statement
	query := statement.SQL.String()

	assertQueryContains(t, query,
		"cp.platform",
		"cp.platform_channel_id",
		"cp.channel_id = $1",
		"cp.platform IN ($2,$3)",
	)
	assertQueryExcludes(t, query,
		"twitch_user_id",
		"kick_user_id",
		"twitch_bot_enabled",
		"kick_bot_enabled",
		"cp.enabled",
		`c."isEnabled"`,
		"vk",
	)
	assertSingleBindingSource(t, query)
	assertQueryVars(t, statement.Vars, "channel-id", platform.PlatformTwitch, platform.PlatformKick)
}

func newDryRunPostgresDB(t *testing.T) *gorm.DB {
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

func assertQueryContains(t *testing.T, query string, fragments ...string) {
	t.Helper()

	for _, fragment := range fragments {
		if !strings.Contains(query, fragment) {
			t.Fatalf("query does not contain %q: %s", fragment, query)
		}
	}
}

func assertQueryExcludes(t *testing.T, query string, fragments ...string) {
	t.Helper()

	for _, fragment := range fragments {
		if strings.Contains(query, fragment) {
			t.Fatalf("query must not contain %q: %s", fragment, query)
		}
	}
}

func assertSingleBindingSource(t *testing.T, query string) {
	t.Helper()

	if got := strings.Count(query, "channel_platforms"); got != 1 {
		t.Fatalf("channel_platforms references = %d, want 1 selected binding source: %s", got, query)
	}
}

func assertQueryVars(t *testing.T, actual []any, expected ...any) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("query vars = %#v, want %#v", actual, expected)
	}
}
