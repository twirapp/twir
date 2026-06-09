package pgx

import (
	"strings"
	"testing"

	"github.com/twirapp/twir/libs/repositories/timers"
)

func TestBuildGetManyQuerySelectsPlatformsBeforeResponses(t *testing.T) {
	query, _, err := buildGetManyQuery(timers.GetManyInput{})
	if err != nil {
		t.Fatalf("build query: %v", err)
	}

	platformsIndex := strings.Index(query, "t.platforms")
	if platformsIndex == -1 {
		t.Fatalf("query does not select t.platforms: %s", query)
	}

	responsesIndex := strings.Index(query, "AS responses")
	if responsesIndex == -1 {
		t.Fatalf("query does not select responses JSON: %s", query)
	}

	if platformsIndex > responsesIndex {
		t.Fatalf("t.platforms must be selected before responses JSON to match GetMany scan order: %s", query)
	}
}
