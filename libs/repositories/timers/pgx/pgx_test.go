package pgx

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/platform"
	channelspgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	"github.com/twirapp/twir/libs/repositories/timers"
)

const integrationDatabaseURLEnv = "TWIR_REPOSITORY_TEST_DATABASE_URL"

func newIntegrationPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	databaseURL := os.Getenv(integrationDatabaseURLEnv)
	if databaseURL == "" {
		t.Skipf("set %s to run PostgreSQL integration tests", integrationDatabaseURLEnv)
	}

	connConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("parse PostgreSQL config: %v", err)
	}
	connConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		platformType, err := conn.LoadType(ctx, "platform")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(platformType)

		platformArrayType, err := conn.LoadType(ctx, "_platform")
		if err != nil {
			return err
		}
		conn.TypeMap().RegisterType(platformArrayType)

		conn.TypeMap().RegisterDefaultPgType(platform.Platform(""), "platform")
		conn.TypeMap().RegisterDefaultPgType([]platform.Platform{}, "_platform")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		t.Fatalf("create PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping PostgreSQL: %v", err)
	}

	return pool
}

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

func TestPgx_Create_defaults_nil_platforms_to_empty_array(t *testing.T) {
	// Given
	ctx := context.Background()
	pool := newIntegrationPool(t, ctx)
	channel, err := channelspgx.New(channelspgx.Opts{PgxPool: pool}).Create(ctx)
	if err != nil {
		t.Fatalf("create test channel: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM channels WHERE id = $1`, channel.ID); err != nil {
			t.Errorf("delete test channel: %v", err)
		}
	})
	repository := New(Opts{Pgx: pool})

	// When
	created, err := repository.Create(ctx, timers.CreateInput{
		ChannelID:       channel.ID.String(),
		Name:            "integration-test-timer",
		Enabled:         true,
		OfflineEnabled:  true,
		OnlineEnabled:   true,
		TimeInterval:    60,
		MessageInterval: 1,
	})

	// Then
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM channels_timers WHERE id = $1`, created.ID); err != nil {
			t.Errorf("delete test timer: %v", err)
		}
	})
	if created.Platforms == nil || len(created.Platforms) != 0 {
		t.Fatalf("platforms = %v, want empty array", created.Platforms)
	}
}
