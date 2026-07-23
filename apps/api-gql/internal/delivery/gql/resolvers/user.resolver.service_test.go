package resolvers

import (
	"context"
	"reflect"
	"strings"
	"testing"

	model "github.com/twirapp/twir/libs/gomodels"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestOwnedDashboardsQueryIncludesZeroBindingLegacyOwners(t *testing.T) {
	t.Parallel()

	const userID = "a8ff9a5f-b450-44eb-ba6c-5047888e3af7"
	statement := ownedDashboardsQuery(
		newResolverDryRunPostgresDB(t),
		context.Background(),
		userID,
	).Find(&[]model.Channels{}).Statement
	query := statement.SQL.String()

	for _, fragment := range []string{
		"EXISTS (",
		"FROM channel_platforms AS cp_owner",
		"cp_owner.channel_id = channels.id",
		"cp_owner.user_id = $1::uuid",
		"NOT EXISTS (",
		"FROM channel_platforms AS cp_existing",
		"cp_existing.channel_id = channels.id",
		"channels.twitch_user_id = $2::uuid",
		"channels.kick_user_id = $3::uuid",
	} {
		if !strings.Contains(query, fragment) {
			t.Fatalf("query does not contain %q: %s", fragment, query)
		}
	}

	if want := []any{userID, userID, userID}; !reflect.DeepEqual(statement.Vars, want) {
		t.Fatalf("query vars = %#v, want %#v", statement.Vars, want)
	}
}

func newResolverDryRunPostgresDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(
		postgres.Open("host=127.0.0.1 user=twir dbname=twir sslmode=disable"),
		&gorm.Config{DisableAutomaticPing: true, DryRun: true},
	)
	if err != nil {
		t.Fatalf("open dry-run database: %v", err)
	}

	return db
}
