package timers

import (
	"context"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestBuildTwitchChannelsQueryUsesJoinedUserPlatformID(t *testing.T) {
	t.Parallel()

	if !strings.Contains(twitchChannelsSelectClause, "users.platform_id AS twitch_platform_id") {
		t.Fatalf("select clause must alias users.platform_id as twitch_platform_id: %s", twitchChannelsSelectClause)
	}

	if !strings.Contains(twitchChannelsJoinClause, "users.id = channels.twitch_user_id") {
		t.Fatalf("join clause must join users through channels.twitch_user_id: %s", twitchChannelsJoinClause)
	}

	if twitchChannelsPlatformIDIsNotNull != "users.platform_id IS NOT NULL" {
		t.Fatalf("not-null clause must target users.platform_id: %s", twitchChannelsPlatformIDIsNotNull)
	}

	db, err := gorm.Open(
		postgres.Open("host=127.0.0.1 user=twir dbname=twir sslmode=disable"),
		&gorm.Config{DisableAutomaticPing: true, DryRun: true},
	)
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}

	statement := buildTwitchChannelsQuery(db, context.Background()).Scan(&[]twitchStreamChannelRow{}).Statement
	sql := statement.SQL.String()

	if !strings.Contains(sql, "twitch_bot_enabled IS TRUE") {
		t.Fatalf("query must filter channels with enabled twitch bot: %s", sql)
	}
	if !strings.Contains(sql, `COALESCE(users.is_banned, false) = false`) {
		t.Fatalf("query must filter banned users: %s", sql)
	}
}

func TestBuildKickChannelsQueryUsesJoinedUserPlatformID(t *testing.T) {
	t.Parallel()

	if strings.Contains(kickChannelsSelectClause, "channels.kick_platform_id") {
		t.Fatalf("select clause must not read channels.kick_platform_id: %s", kickChannelsSelectClause)
	}

	if !strings.Contains(kickChannelsSelectClause, "users.platform_id AS kick_platform_id") {
		t.Fatalf("select clause must alias users.platform_id as kick_platform_id: %s", kickChannelsSelectClause)
	}

	if !strings.Contains(kickChannelsJoinClause, "users.id = channels.kick_user_id") {
		t.Fatalf("join clause must join users through channels.kick_user_id: %s", kickChannelsJoinClause)
	}

	if kickChannelsPlatformIDIsNotNull != "users.platform_id IS NOT NULL" {
		t.Fatalf("not-null clause must target users.platform_id: %s", kickChannelsPlatformIDIsNotNull)
	}
}
