package postgres

import (
	"strings"
	"testing"
)

func TestBuildGetByUserAndChannelIDQueryComparesIDs(t *testing.T) {
	query := buildGetByUserAndChannelIDQuery()

	if !strings.Contains(query, `user_id = $1`) {
		t.Fatalf("query must compare user_id: %s", query)
	}

	if !strings.Contains(query, `channel_id = $2`) {
		t.Fatalf("query must compare channel_id: %s", query)
	}
}
