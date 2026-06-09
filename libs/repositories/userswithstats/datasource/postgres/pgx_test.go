package postgres

import (
	"strings"
	"testing"
)

func TestGetByUserAndChannelIDQueryCastsStatsIDsToText(t *testing.T) {
	if !strings.Contains(getByUserAndChannelIDQuery, `u.id::text = $1::text`) {
		t.Fatalf("query must compare users id via text cast: %s", getByUserAndChannelIDQuery)
	}

	if !strings.Contains(getByUserAndChannelIDQuery, `us."userId"::text = $1::text`) {
		t.Fatalf("query must compare users_stats userId via text cast: %s", getByUserAndChannelIDQuery)
	}

	if !strings.Contains(getByUserAndChannelIDQuery, `us."channelId"::text = $2::text`) {
		t.Fatalf("query must compare users_stats channelId via text cast: %s", getByUserAndChannelIDQuery)
	}
}
