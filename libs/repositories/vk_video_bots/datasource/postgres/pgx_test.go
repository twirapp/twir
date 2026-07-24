package postgres

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	entity "github.com/twirapp/twir/libs/entities/vk_video_bot"
	vkvideobots "github.com/twirapp/twir/libs/repositories/vk_video_bots"
)

func TestPgxUpsertAcquiresSingletonLockAndReturnsStoredBot(t *testing.T) {
	// Given
	want := testVKVideoBot()
	executor := &vkVideoBotExecutor{
		queryFn: func(_ context.Context, query string, args ...any) (pgx.Rows, error) {
			lowerQuery := strings.ToLower(query)
			for _, fragment := range []string{
				"pg_advisory_xact_lock",
				"insert into vk_video_bots",
				"on conflict (singleton)",
				"encrypted_access_token",
				"encrypted_refresh_token",
				"returning",
			} {
				if !strings.Contains(lowerQuery, fragment) {
					t.Fatalf("upsert query missing %q: %s", fragment, query)
				}
			}
			if len(args) != 6 || args[0] != want.EncryptedAccessToken || args[1] != want.EncryptedRefreshToken || args[5] != want.VKUserID {
				t.Fatalf("upsert args = %#v, want encrypted token data and VK user ID", args)
			}

			return &vkVideoBotRows{bot: &want}, nil
		},
	}
	repository := &Pgx{pool: executor}

	// When
	got, err := repository.Upsert(context.Background(), vkvideobots.UpsertInput{
		EncryptedAccessToken:  want.EncryptedAccessToken,
		EncryptedRefreshToken: want.EncryptedRefreshToken,
		Scopes:                want.Scopes,
		ExpiresIn:             want.ExpiresIn,
		ObtainmentTimestamp:   want.ObtainmentTimestamp,
		VKUserID:              want.VKUserID,
	})

	// Then
	if err != nil {
		t.Fatalf("upsert VK Video bot: %v", err)
	}
	requireVKVideoBotEqual(t, got, want)
}

func TestPgxGetAcquiresSingletonLockBeforeReadingBot(t *testing.T) {
	// Given
	want := testVKVideoBot()
	executor := &vkVideoBotExecutor{
		queryFn: func(_ context.Context, query string, _ ...any) (pgx.Rows, error) {
			lowerQuery := strings.ToLower(query)
			lockIndex := strings.Index(lowerQuery, "pg_advisory_xact_lock")
			readIndex := strings.LastIndex(lowerQuery, "select")
			if lockIndex == -1 || readIndex == -1 || lockIndex > readIndex || !strings.Contains(lowerQuery, "where singleton") {
				t.Fatalf("get query does not lock and select the singleton: %s", query)
			}

			return &vkVideoBotRows{bot: &want}, nil
		},
	}
	repository := &Pgx{pool: executor}

	// When
	got, err := repository.Get(context.Background())

	// Then
	if err != nil {
		t.Fatalf("get VK Video bot: %v", err)
	}
	requireVKVideoBotEqual(t, got, want)
}

func requireVKVideoBotEqual(t *testing.T, got, want entity.VKVideoBot) {
	t.Helper()
	if got.ID != want.ID || got.EncryptedAccessToken != want.EncryptedAccessToken ||
		got.EncryptedRefreshToken != want.EncryptedRefreshToken || got.ExpiresIn != want.ExpiresIn ||
		got.VKUserID != want.VKUserID || !slices.Equal(got.Scopes, want.Scopes) ||
		!got.ObtainmentTimestamp.Equal(want.ObtainmentTimestamp) || !got.CreatedAt.Equal(want.CreatedAt) ||
		!got.UpdatedAt.Equal(want.UpdatedAt) {
		t.Fatalf("bot = %#v, want %#v", got, want)
	}
}

func testVKVideoBot() entity.VKVideoBot {
	createdAt := time.Date(2026, 7, 24, 0, 0, 0, 0, time.UTC)
	return entity.VKVideoBot{
		ID:                    uuid.New(),
		EncryptedAccessToken:  "encrypted-access",
		EncryptedRefreshToken: "encrypted-refresh",
		Scopes:                []string{"chat"},
		ExpiresIn:             3600,
		ObtainmentTimestamp:   createdAt,
		VKUserID:              uuid.New(),
		CreatedAt:             createdAt,
		UpdatedAt:             createdAt,
	}
}

type vkVideoBotExecutor struct {
	queryFn func(context.Context, string, ...any) (pgx.Rows, error)
}

func (e *vkVideoBotExecutor) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return e.queryFn(ctx, query, args...)
}

func (*vkVideoBotExecutor) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}

type vkVideoBotRows struct {
	bot      *entity.VKVideoBot
	returned bool
}

func (*vkVideoBotRows) Close() {}

func (*vkVideoBotRows) Err() error { return nil }

func (*vkVideoBotRows) CommandTag() pgconn.CommandTag { return pgconn.NewCommandTag("") }

func (*vkVideoBotRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{
		{Name: "id"},
		{Name: "encrypted_access_token"},
		{Name: "encrypted_refresh_token"},
		{Name: "scopes"},
		{Name: "expires_in"},
		{Name: "obtainment_timestamp"},
		{Name: "vk_user_id"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}
}

func (r *vkVideoBotRows) Next() bool {
	if r.bot == nil || r.returned {
		return false
	}
	r.returned = true
	return true
}

func (r *vkVideoBotRows) Scan(dest ...any) error {
	if r.bot == nil {
		return pgx.ErrNoRows
	}
	if len(dest) != 9 {
		return fmt.Errorf("scan destinations = %d, want 9", len(dest))
	}
	*dest[0].(*uuid.UUID) = r.bot.ID
	*dest[1].(*string) = r.bot.EncryptedAccessToken
	*dest[2].(*string) = r.bot.EncryptedRefreshToken
	*dest[3].(*[]string) = append([]string(nil), r.bot.Scopes...)
	*dest[4].(*int) = r.bot.ExpiresIn
	*dest[5].(*time.Time) = r.bot.ObtainmentTimestamp
	*dest[6].(*uuid.UUID) = r.bot.VKUserID
	*dest[7].(*time.Time) = r.bot.CreatedAt
	*dest[8].(*time.Time) = r.bot.UpdatedAt
	return nil
}

func (*vkVideoBotRows) Values() ([]any, error) { return nil, nil }

func (*vkVideoBotRows) RawValues() [][]byte { return nil }

func (*vkVideoBotRows) Conn() *pgx.Conn { return nil }
