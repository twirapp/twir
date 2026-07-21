package pgx

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/platform"
)

const channelPlatformsMigrationTestDatabaseURLEnv = "TWIR_CHANNELS_MULTIPLATFORM_MIGRATION_TEST_DATABASE_URL"

func TestChannelPlatformsMigrationContainsSafeBackfills(t *testing.T) {
	migration := channelPlatformsMigrationUpSQL(t)

	if got := strings.Count(migration, "INSERT INTO channel_platforms"); got != 2 {
		t.Fatalf("channel_platforms backfill statements = %d, want 2", got)
	}

	if got := strings.Count(migration, "ON CONFLICT DO NOTHING"); got != 2 {
		t.Fatalf("channel_platforms conflict-safe inserts = %d, want 2", got)
	}

	for _, fragment := range []string{
		"'twitch'",
		"'kick'",
		"kick_bot_id",
		"RAISE EXCEPTION",
	} {
		if !strings.Contains(migration, fragment) {
			t.Fatalf("migration does not validate safe backfill fragment %q", fragment)
		}
	}
}

func TestChannelPlatformsMigrationBackfillsLegacyBindings(t *testing.T) {
	databaseURL := os.Getenv(channelPlatformsMigrationTestDatabaseURLEnv)
	if databaseURL == "" {
		t.Skipf("set %s to run PostgreSQL migration tests", channelPlatformsMigrationTestDatabaseURLEnv)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping PostgreSQL: %v", err)
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("acquire PostgreSQL connection: %v", err)
	}
	t.Cleanup(conn.Release)

	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatalf("begin migration test transaction: %v", err)
	}
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
	})

	setupMigrationTestSchema(t, ctx, tx)
	fixtures := insertLegacyChannelFixtures(t, ctx, tx)

	if _, err := tx.Exec(ctx, channelPlatformsMigrationUpSQL(t), pgx.QueryExecModeSimpleProtocol); err != nil {
		t.Fatalf("run channel_platforms migration: %v", err)
	}

	for _, fixture := range fixtures {
		fixture := fixture
		t.Run(fixture.name, func(t *testing.T) {
			for _, want := range fixture.bindings {
				assertBackfilledBinding(t, ctx, tx, fixture.channelID, want)
			}
		})
	}

	var bindingCount int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM channel_platforms`).Scan(&bindingCount); err != nil {
		t.Fatalf("count backfilled bindings: %v", err)
	}
	if bindingCount != 4 {
		t.Fatalf("backfilled binding count = %d, want 4", bindingCount)
	}
}

type legacyChannelFixture struct {
	name      string
	channelID uuid.UUID
	bindings  []backfilledBinding
}

type backfilledBinding struct {
	platform          platform.Platform
	userID            uuid.UUID
	platformChannelID string
	enabled           bool
	kickBotID         *uuid.UUID
	kickBotUserID     *uuid.UUID
}

func setupMigrationTestSchema(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()

	for _, query := range []string{
		`SET LOCAL search_path TO pg_temp`,
		`CREATE TEMP TABLE users (
			id UUID PRIMARY KEY,
			platform TEXT NOT NULL,
			platform_id TEXT NOT NULL
		)`,
		`CREATE TEMP TABLE kick_bots (
			id UUID PRIMARY KEY,
			kick_user_id UUID NOT NULL REFERENCES users(id)
		)`,
		`CREATE TEMP TABLE channels (
			id UUID PRIMARY KEY,
			twitch_user_id UUID REFERENCES users(id),
			twitch_bot_enabled BOOLEAN NOT NULL DEFAULT false,
			kick_user_id UUID REFERENCES users(id),
			kick_bot_enabled BOOLEAN NOT NULL DEFAULT false,
			"botId" TEXT NOT NULL DEFAULT '',
			"isBotMod" BOOLEAN NOT NULL DEFAULT false,
			"isTwitchBanned" BOOLEAN NOT NULL DEFAULT false,
			kick_bot_id UUID REFERENCES kick_bots(id)
		)`,
		`CREATE FUNCTION uuidv7()
			RETURNS UUID
			LANGUAGE SQL
			VOLATILE
			AS $$
				SELECT md5(random()::text || clock_timestamp()::text)::uuid
			$$`,
	} {
		if _, err := tx.Exec(ctx, query); err != nil {
			t.Fatalf("set up migration test schema: %v", err)
		}
	}
}

func insertLegacyChannelFixtures(t *testing.T, ctx context.Context, tx pgx.Tx) []legacyChannelFixture {
	t.Helper()

	twitchOnlyUserID := insertLegacyUser(t, ctx, tx, platform.PlatformTwitch, "twitch-only-channel")
	kickOnlyUserID := insertLegacyUser(t, ctx, tx, platform.PlatformKick, "kick-only-channel")
	kickOnlyBotUserID := insertLegacyUser(t, ctx, tx, platform.PlatformKick, "kick-only-bot")
	dualTwitchUserID := insertLegacyUser(t, ctx, tx, platform.PlatformTwitch, "dual-twitch-channel")
	dualKickUserID := insertLegacyUser(t, ctx, tx, platform.PlatformKick, "dual-kick-channel")
	dualKickBotUserID := insertLegacyUser(t, ctx, tx, platform.PlatformKick, "dual-kick-bot")

	kickOnlyBotID := insertLegacyKickBot(t, ctx, tx, kickOnlyBotUserID)
	dualKickBotID := insertLegacyKickBot(t, ctx, tx, dualKickBotUserID)

	twitchOnlyChannelID := insertLegacyChannel(
		t,
		ctx,
		tx,
		&twitchOnlyUserID,
		true,
		nil,
		false,
		"twitch-only-bot",
		true,
		false,
		nil,
	)
	kickOnlyChannelID := insertLegacyChannel(
		t,
		ctx,
		tx,
		nil,
		false,
		&kickOnlyUserID,
		false,
		"",
		false,
		false,
		&kickOnlyBotID,
	)
	dualChannelID := insertLegacyChannel(
		t,
		ctx,
		tx,
		&dualTwitchUserID,
		false,
		&dualKickUserID,
		true,
		"dual-twitch-bot",
		true,
		false,
		&dualKickBotID,
	)

	return []legacyChannelFixture{
		{
			name:      "Twitch-only",
			channelID: twitchOnlyChannelID,
			bindings: []backfilledBinding{{
				platform:          platform.PlatformTwitch,
				userID:            twitchOnlyUserID,
				platformChannelID: "twitch-only-channel",
				enabled:           true,
			}},
		},
		{
			name:      "Kick-only",
			channelID: kickOnlyChannelID,
			bindings: []backfilledBinding{{
				platform:          platform.PlatformKick,
				userID:            kickOnlyUserID,
				platformChannelID: "kick-only-channel",
				enabled:           false,
				kickBotID:         &kickOnlyBotID,
				kickBotUserID:     &kickOnlyBotUserID,
			}},
		},
		{
			name:      "dual-bound",
			channelID: dualChannelID,
			bindings: []backfilledBinding{
				{
					platform:          platform.PlatformTwitch,
					userID:            dualTwitchUserID,
					platformChannelID: "dual-twitch-channel",
					enabled:           false,
				},
				{
					platform:          platform.PlatformKick,
					userID:            dualKickUserID,
					platformChannelID: "dual-kick-channel",
					enabled:           true,
					kickBotID:         &dualKickBotID,
					kickBotUserID:     &dualKickBotUserID,
				},
			},
		},
	}
}

func insertLegacyUser(
	t *testing.T,
	ctx context.Context,
	tx pgx.Tx,
	p platform.Platform,
	platformID string,
) uuid.UUID {
	t.Helper()

	userID := uuid.New()
	if _, err := tx.Exec(
		ctx,
		`INSERT INTO users (id, platform, platform_id) VALUES ($1, $2, $3)`,
		userID,
		p,
		platformID,
	); err != nil {
		t.Fatalf("insert legacy user: %v", err)
	}

	return userID
}

func insertLegacyKickBot(t *testing.T, ctx context.Context, tx pgx.Tx, userID uuid.UUID) uuid.UUID {
	t.Helper()

	botID := uuid.New()
	if _, err := tx.Exec(
		ctx,
		`INSERT INTO kick_bots (id, kick_user_id) VALUES ($1, $2)`,
		botID,
		userID,
	); err != nil {
		t.Fatalf("insert legacy Kick bot: %v", err)
	}

	return botID
}

func insertLegacyChannel(
	t *testing.T,
	ctx context.Context,
	tx pgx.Tx,
	twitchUserID *uuid.UUID,
	twitchBotEnabled bool,
	kickUserID *uuid.UUID,
	kickBotEnabled bool,
	botID string,
	isBotMod bool,
	isTwitchBanned bool,
	kickBotID *uuid.UUID,
) uuid.UUID {
	t.Helper()

	channelID := uuid.New()
	if _, err := tx.Exec(
		ctx,
		`INSERT INTO channels (
			id,
			twitch_user_id,
			twitch_bot_enabled,
			kick_user_id,
			kick_bot_enabled,
			"botId",
			"isBotMod",
			"isTwitchBanned",
			kick_bot_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		channelID,
		twitchUserID,
		twitchBotEnabled,
		kickUserID,
		kickBotEnabled,
		botID,
		isBotMod,
		isTwitchBanned,
		kickBotID,
	); err != nil {
		t.Fatalf("insert legacy channel: %v", err)
	}

	return channelID
}

func assertBackfilledBinding(
	t *testing.T,
	ctx context.Context,
	tx pgx.Tx,
	channelID uuid.UUID,
	want backfilledBinding,
) {
	t.Helper()

	var got struct {
		UserID            uuid.UUID
		PlatformChannelID string
		Enabled           bool
		BotUserID         *uuid.UUID
		BotConfig         []byte
	}
	if err := tx.QueryRow(
		ctx,
		`SELECT user_id, platform_channel_id, enabled, bot_user_id, bot_config
		 FROM channel_platforms
		 WHERE channel_id = $1 AND platform = $2`,
		channelID,
		want.platform,
	).Scan(
		&got.UserID,
		&got.PlatformChannelID,
		&got.Enabled,
		&got.BotUserID,
		&got.BotConfig,
	); err != nil {
		t.Fatalf("get %s binding: %v", want.platform, err)
	}

	if got.UserID != want.userID {
		t.Fatalf("%s binding user ID = %s, want %s", want.platform, got.UserID, want.userID)
	}
	if got.PlatformChannelID != want.platformChannelID {
		t.Fatalf("%s binding platform channel ID = %q, want %q", want.platform, got.PlatformChannelID, want.platformChannelID)
	}
	if got.Enabled != want.enabled {
		t.Fatalf("%s binding enabled = %t, want %t", want.platform, got.Enabled, want.enabled)
	}

	if want.kickBotID == nil {
		return
	}

	if got.BotUserID == nil || *got.BotUserID != *want.kickBotUserID {
		t.Fatalf("Kick binding bot user ID = %v, want %s", got.BotUserID, want.kickBotUserID)
	}

	var botConfig map[string]string
	if err := json.Unmarshal(got.BotConfig, &botConfig); err != nil {
		t.Fatalf("unmarshal Kick bot config: %v", err)
	}
	if gotKickBotID := botConfig["kick_bot_id"]; gotKickBotID != want.kickBotID.String() {
		t.Fatalf("Kick bot config ID = %q, want %q", gotKickBotID, want.kickBotID.String())
	}
}

func channelPlatformsMigrationUpSQL(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate migration test source")
	}

	migrationPath := filepath.Join(
		filepath.Dir(thisFile),
		"../../../migrations/postgres/20260721120000_channel_platforms.sql",
	)
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("read channel_platforms migration: %v", err)
	}

	up, _, found := strings.Cut(string(migration), "-- +goose Down")
	if !found {
		t.Fatal("channel_platforms migration does not contain a Goose down section")
	}

	lines := make([]string, 0)
	for _, line := range strings.Split(up, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "-- +goose ") {
			continue
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
