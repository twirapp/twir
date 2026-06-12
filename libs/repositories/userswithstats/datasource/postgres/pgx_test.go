package postgres

import (
	"strings"
	"testing"
)

func TestGetByUserAndChannelIDQueryCastsStatsIDsToText(t *testing.T) {
	if !strings.Contains(getByUserAndChannelIDQuery, `u.id = $1`) {
		t.Fatalf("query must compare users id: %s", getByUserAndChannelIDQuery)
	}

	if !strings.Contains(getByUserAndChannelIDQuery, `us.user_id = $1`) {
		t.Fatalf("query must compare users_stats user_id: %s", getByUserAndChannelIDQuery)
	}

	if !strings.Contains(getByUserAndChannelIDQuery, `us.channel_id = $2`) {
		t.Fatalf("query must compare users_stats channel_id: %s", getByUserAndChannelIDQuery)
	}
}
