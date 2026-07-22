package pgx

import (
	"strings"
	"testing"
)

func TestPatchQueryUpdatesOnlyBindingState(t *testing.T) {
	query := strings.ToLower(patchQuery)

	for _, want := range []string{
		"update channel_platforms",
		"enabled = coalesce($2::boolean, enabled)",
		"coalesce(bot_config, '{}'::jsonb) || coalesce($3::jsonb, '{}'::jsonb)",
		"where id = $1",
	} {
		if !strings.Contains(query, want) {
			t.Fatalf("patch query missing %q: %s", want, patchQuery)
		}
	}

	for _, forbidden := range []string{"user_id =", "platform_channel_id =", "bot_user_id ="} {
		if strings.Contains(query, forbidden) {
			t.Fatalf("patch query overwrites unrelated binding field %q: %s", forbidden, patchQuery)
		}
	}
}
