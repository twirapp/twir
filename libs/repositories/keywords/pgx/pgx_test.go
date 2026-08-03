package pgx

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/platform"
	channelspgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	"github.com/twirapp/twir/libs/repositories/keywords"
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
	repository := New(Opts{PgxPool: pool})

	// When
	created, err := repository.Create(ctx, keywords.CreateInput{
		ChannelID: channel.ID,
		Text:      "integration-test-keyword",
		Enabled:   true,
	})

	// Then
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	t.Cleanup(func() {
		if err := repository.Delete(ctx, created.ID); err != nil {
			t.Errorf("delete test keyword: %v", err)
		}
	})
	if created.Platforms == nil || len(created.Platforms) != 0 {
		t.Fatalf("platforms = %v, want empty array", created.Platforms)
	}
}

func TestPgx_Update_preserves_platforms_when_input_omits_them(t *testing.T) {
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
	repository := New(Opts{PgxPool: pool})
	created, err := repository.Create(ctx, keywords.CreateInput{
		ChannelID: channel.ID,
		Text:      "integration-test-keyword",
		Enabled:   true,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	t.Cleanup(func() {
		if err := repository.Delete(ctx, created.ID); err != nil {
			t.Errorf("delete test keyword: %v", err)
		}
	})
	newText := "integration-test-keyword-updated"

	// When
	updated, err := repository.Update(ctx, created.ID, keywords.UpdateInput{Text: &newText})

	// Then
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Platforms == nil || len(updated.Platforms) != 0 {
		t.Fatalf("platforms = %v, want preserved empty array", updated.Platforms)
	}
	if updated.Text != newText {
		t.Fatalf("text = %q, want %q", updated.Text, newText)
	}
}
