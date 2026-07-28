package timers

import (
	"context"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBuildTwitchChannelsQueryUsesTwitchBinding(t *testing.T) {
	t.Parallel()

	db := newDryRunPostgresDB(t)
	statement := buildTwitchChannelsQuery(db, context.Background()).Scan(&[]twitchStreamChannelRow{}).Statement
	sql := statement.SQL.String()

	assertQueryContains(t, sql,
		"cp.channel_id AS id",
		"cp.platform_channel_id AS twitch_platform_id",
		"JOIN channels c ON c.id = cp.channel_id",
		"JOIN users u ON u.id = cp.user_id AND u.platform = 'twitch'",
		"cp.platform = $1",
		"cp.enabled = $2",
		"COALESCE(u.is_banned, false) = false",
	)
	assertQueryExcludes(t, sql, "twitch_user_id", "twitch_bot_enabled", "users.platform_id")
	assertSingleBindingSource(t, sql)
}

func TestBuildKickChannelsQueryUsesKickBindingAndRetainsGlobalEnablementEligibility(t *testing.T) {
	t.Parallel()

	db := newDryRunPostgresDB(t)
	statement := buildKickChannelsQuery(db, context.Background()).Scan(&[]kickChannelRow{}).Statement
	sql := statement.SQL.String()

	assertQueryContains(t, sql,
		"cp.channel_id AS id",
		"cp.platform_channel_id AS kick_platform_id",
		"JOIN channels c ON c.id = cp.channel_id",
		"JOIN users u ON u.id = cp.user_id AND u.platform = 'kick'",
		"cp.platform = $1",
		`c."isEnabled" = $2`,
		"COALESCE(u.is_banned, false) = false",
	)
	assertQueryExcludes(t, sql, "kick_user_id", "kick_bot_enabled", "users.platform_id", "cp.enabled")
	assertSingleBindingSource(t, sql)
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
