package emotes_store

import (
	"context"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBuildStartupChannelsQueryUsesExplicitNormalizedBindings(t *testing.T) {
	t.Parallel()

	db := newDryRunPostgresDB(t)
	statement := buildStartupChannelsQuery(db, context.Background()).Scan(&[]startupChannelData{}).Statement
	query := statement.SQL.String()

	assertQueryContains(t, query,
		"FROM channel_platforms AS cp",
		"cp.channel_id AS id",
		"cp.platform",
		"cp.platform_channel_id AS platform_id",
		"JOIN channels c ON c.id = cp.channel_id",
		`c."isEnabled" = $1`,
		"cp.platform IN ($2,$3)",
	)
	assertQueryExcludes(t, query,
		"twitch_user_id",
		"kick_user_id",
		"twitch_platform_id",
		"kick_platform_id",
		"users.platform_id",
		"cp.enabled",
	)
	assertSingleBindingSource(t, query)
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
