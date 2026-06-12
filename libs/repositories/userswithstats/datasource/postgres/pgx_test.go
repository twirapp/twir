package postgres

import (
	"strings"
	"testing"
)

func TestGetByUserAndChannelIDQueryCastsStatsIDsToText(t *testing.T) {
	if !strings.Contains(getByUserAndChannelIDQuery, `u.id = $1::uuid`) {
		t.Fatalf("query must compare users id via uuid: %s", getByUserAndChannelIDQuery)
	}

	if !strings.Contains(getByUserAndChannelIDQuery, `us.user_id = $1::uuid`) {
		t.Fatalf("query must compare users_stats user_id via uuid: %s", getByUserAndChannelIDQuery)
	}

	if !strings.Contains(getByUserAndChannelIDQuery, `us.channel_id = $2::uuid`) {
		t.Fatalf("query must compare users_stats channel_id via uuid: %s", getByUserAndChannelIDQuery)
	}
}
