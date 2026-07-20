package pgx

import (
	"strings"
	"testing"

	"github.com/twirapp/twir/libs/repositories/users_with_channel"
)

func TestBuildGetManyQuerySkipsBadgesJoinWithoutBadgeFilter(t *testing.T) {
	query, _, err := buildGetManyQuery(users_with_channel.GetManyInput{PerPage: 10})
	if err != nil {
		t.Fatalf("build query: %v", err)
	}
	if strings.Contains(query, "badges_users") {
		t.Fatalf("unexpected badges join in query: %s", query)
	}
}

func TestBuildGetManyCountQueryWithoutJoinFiltersCountsUsersDirectly(t *testing.T) {
	query, args, err := buildGetManyCountQuery(users_with_channel.GetManyInput{})
	if err != nil {
		t.Fatalf("build query: %v", err)
	}
	if query != "SELECT COUNT(*) FROM users u" {
		t.Fatalf("unexpected count query: %s", query)
	}
	if len(args) != 0 {
		t.Fatalf("expected no query arguments, got %d", len(args))
	}
}
