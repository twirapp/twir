package timers

import (
	"strings"
	"testing"
)

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
