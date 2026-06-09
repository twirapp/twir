package postgres

import (
	"strings"
	"testing"
)

func TestBuildGetByUserAndChannelIDQueryCastsIDsToText(t *testing.T) {
	query := buildGetByUserAndChannelIDQuery()

	if !strings.Contains(query, `"userId"::text = $1::text`) {
		t.Fatalf("query must compare userId via text cast: %s", query)
	}

	if !strings.Contains(query, `"channelId"::text = $2::text`) {
		t.Fatalf("query must compare channelId via text cast: %s", query)
	}
}
