package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/channel_platforms"
	"github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	"github.com/twirapp/twir/libs/repositories/channels"
	channelspgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	"github.com/twirapp/twir/libs/repositories/users"
	userspgx "github.com/twirapp/twir/libs/repositories/users/pgx"
)

func TestChannelPlatformNil(t *testing.T) {
	if !model.Nil.IsNil() {
		t.Fatal("Nil must report itself as nil")
	}

	if (model.ChannelPlatform{}).IsNil() {
		t.Fatal("a populated-value model must not report itself as nil")
	}
}

func TestRequireUniqueViolationMatchesPgxV5Error(t *testing.T) {
	const childEnv = "TWIR_CHANNEL_PLATFORMS_UNIQUE_VIOLATION_CHILD"

	if os.Getenv(childEnv) == "1" {
		requireUniqueViolation(
			t,
			fmt.Errorf("wrapped: %w", &pgconn.PgError{Code: "23505"}),
		)
		return
	}

	command := exec.Command(os.Args[0], "-test.run=^TestRequireUniqueViolationMatchesPgxV5Error$")
	command.Env = append(os.Environ(), childEnv+"=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("requireUniqueViolation rejected a pgx v5 unique violation: %v\n%s", err, output)
	}
}

func TestRepositoryIntegration(t *testing.T) {
	ctx := context.Background()
	repository, pool := newIntegrationRepository(t, ctx)

	firstUserID := createTestUser(t, ctx, pool)
	secondUserID := createTestUser(t, ctx, pool)
	firstChannelID := createTestChannel(t, ctx, pool, firstUserID)
	secondChannelID := createTestChannel(t, ctx, pool, secondUserID)

	t.Cleanup(func() {
		for _, channelID := range []uuid.UUID{firstChannelID, secondChannelID} {
			if _, err := pool.Exec(ctx, `DELETE FROM channel_platforms WHERE channel_id = $1`, channelID); err != nil {
				t.Errorf("delete channel platform bindings for cleanup: %v", err)
			}
		}
	})

	created, err := repository.Create(ctx, channel_platforms.CreateInput{
		ChannelID:         firstChannelID,
		Platform:          platform.PlatformVKVideoLive,
		UserID:            firstUserID,
		PlatformChannelID: "vk-channel-1",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{"mode":"primary"}`),
	})
	if err != nil {
		t.Fatalf("create binding: %v", err)
	}

	if created.ChannelID != firstChannelID || created.Platform != platform.PlatformVKVideoLive || created.UserID != firstUserID || created.PlatformChannelID != "vk-channel-1" || !created.Enabled {
		t.Fatalf("created binding = %+v", created)
	}

	t.Run("GetByChannelAndPlatform", func(t *testing.T) {
		binding, err := repository.GetByChannelAndPlatform(ctx, firstChannelID, platform.PlatformVKVideoLive)
		if err != nil {
			t.Fatalf("get binding by channel and platform: %v", err)
		}
		if binding.ID != created.ID {
			t.Fatalf("binding ID = %s, want %s", binding.ID, created.ID)
		}
	})

	t.Run("GetByPlatformChannelID", func(t *testing.T) {
		binding, err := repository.GetByPlatformChannelID(ctx, platform.PlatformVKVideoLive, "vk-channel-1")
		if err != nil {
			t.Fatalf("get binding by provider channel ID: %v", err)
		}
		if binding.ID != created.ID {
			t.Fatalf("binding ID = %s, want %s", binding.ID, created.ID)
		}
	})

	if _, err := repository.Create(ctx, channel_platforms.CreateInput{
		ChannelID:         firstChannelID,
		Platform:          platform.PlatformVKVideoLive,
		UserID:            secondUserID,
		PlatformChannelID: "vk-channel-2",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{}`),
	}); err == nil {
		t.Fatal("creating a duplicate channel/platform binding succeeded")
	} else {
		requireUniqueViolation(t, err)
	}

	if _, err := repository.Create(ctx, channel_platforms.CreateInput{
		ChannelID:         firstChannelID,
		Platform:          platform.PlatformTwitch,
		UserID:            firstUserID,
		PlatformChannelID: "twitch-channel-1",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{}`),
	}); err != nil {
		t.Fatalf("create second binding: %v", err)
	}

	t.Run("ListByChannelID", func(t *testing.T) {
		bindings, err := repository.ListByChannelID(ctx, firstChannelID)
		if err != nil {
			t.Fatalf("list bindings: %v", err)
		}
		if len(bindings) != 2 {
			t.Fatalf("binding count = %d, want 2", len(bindings))
		}
	})

	t.Run("Update", func(t *testing.T) {
		updated, err := repository.Update(ctx, created.ID, channel_platforms.UpdateInput{
			UserID:            firstUserID,
			PlatformChannelID: "vk-channel-1-updated",
			Enabled:           false,
			BotUserID:         &secondUserID,
			BotConfig:         json.RawMessage(`{"mode":"updated"}`),
		})
		if err != nil {
			t.Fatalf("update binding: %v", err)
		}
		if updated.Enabled || updated.PlatformChannelID != "vk-channel-1-updated" || updated.BotUserID == nil || *updated.BotUserID != secondUserID {
			t.Fatalf("updated binding = %+v", updated)
		}
	})

	if _, err := repository.Create(ctx, channel_platforms.CreateInput{
		ChannelID:         secondChannelID,
		Platform:          platform.PlatformVKVideoLive,
		UserID:            secondUserID,
		PlatformChannelID: "vk-channel-1-updated",
		Enabled:           true,
		BotConfig:         json.RawMessage(`{}`),
	}); err == nil {
		t.Fatal("linking a provider account to a second Twir channel succeeded")
	} else {
		requireUniqueViolation(t, err)
	}

	t.Run("Delete", func(t *testing.T) {
		if err := repository.Delete(ctx, created.ID); err != nil {
			t.Fatalf("delete binding: %v", err)
		}

		binding, err := repository.GetByChannelAndPlatform(ctx, firstChannelID, platform.PlatformVKVideoLive)
		if !errors.Is(err, channel_platforms.ErrNotFound) {
			t.Fatalf("get deleted binding error = %v, want ErrNotFound", err)
		}
		if !binding.IsNil() {
			t.Fatalf("get deleted binding = %+v, want Nil", binding)
		}
	})
}

func newIntegrationRepository(t *testing.T, ctx context.Context) (*Pgx, *pgxpool.Pool) {
	t.Helper()

	databaseURL := os.Getenv("TWIR_CHANNEL_PLATFORMS_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("set TWIR_CHANNEL_PLATFORMS_TEST_DATABASE_URL to run PostgreSQL integration tests")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create PostgreSQL pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping PostgreSQL: %v", err)
	}

	return New(Opts{PgxPool: pool}), pool
}

func createTestUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	user, err := userspgx.New(userspgx.Opts{PgxPool: pool}).Create(ctx, users.CreateInput{
		Platform:          platform.PlatformTwitch,
		PlatformID:        uuid.NewString(),
		IsBotAdmin:        false,
		IsBanned:          false,
		HideOnLandingPage: true,
		Login:             "channel-platform-test-" + uuid.NewString(),
		DisplayName:       "Channel Platform Test",
		Avatar:            "",
	})
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}

	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, user.ID); err != nil {
			t.Errorf("delete test user: %v", err)
		}
	})

	return user.ID
}

func createTestChannel(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) uuid.UUID {
	t.Helper()

	channel, err := channelspgx.New(channelspgx.Opts{PgxPool: pool}).Create(ctx, channels.CreateInput{
		TwitchUserID:     &userID,
		TwitchBotEnabled: true,
		KickBotEnabled:   false,
		BotID:            "channel-platform-test-" + uuid.NewString(),
	})
	if err != nil {
		t.Fatalf("create test channel: %v", err)
	}

	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, `DELETE FROM channels WHERE id = $1`, channel.ID); err != nil {
			t.Errorf("delete test channel: %v", err)
		}
	})

	return channel.ID
}

func requireUniqueViolation(t *testing.T, err error) {
	t.Helper()

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		t.Fatalf("error = %v, want unique-constraint violation", err)
	}
}
