package pgx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/twirapp/twir/libs/entities/platform"
	channelplatforms "github.com/twirapp/twir/libs/repositories/channel_platforms"
	"github.com/twirapp/twir/libs/repositories/channel_platforms/model"
)

type patchExecutor struct {
	queryCalls int
	queryFn    func(context.Context, string, ...any) (pgx.Rows, error)
}

func (e *patchExecutor) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	e.queryCalls++
	return e.queryFn(ctx, query, args...)
}

func (e *patchExecutor) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}

type patchRows struct {
	binding  *model.ChannelPlatform
	returned bool
}

func (r *patchRows) Close() {}

func (r *patchRows) Err() error {
	return nil
}

func (r *patchRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("")
}

func (r *patchRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{
		{Name: "id"},
		{Name: "channel_id"},
		{Name: "platform"},
		{Name: "user_id"},
		{Name: "platform_channel_id"},
		{Name: "enabled"},
		{Name: "bot_user_id"},
		{Name: "bot_config"},
		{Name: "created_at"},
		{Name: "updated_at"},
	}
}

func (r *patchRows) Next() bool {
	if r.binding == nil || r.returned {
		return false
	}

	r.returned = true
	return true
}

func (r *patchRows) Scan(dest ...any) error {
	if r.binding == nil {
		return pgx.ErrNoRows
	}
	if len(dest) != 10 {
		return fmt.Errorf("scan destinations = %d, want 10", len(dest))
	}

	*dest[0].(*uuid.UUID) = r.binding.ID
	*dest[1].(*uuid.UUID) = r.binding.ChannelID
	*dest[2].(*platform.Platform) = r.binding.Platform
	*dest[3].(*uuid.UUID) = r.binding.UserID
	*dest[4].(*string) = r.binding.PlatformChannelID
	*dest[5].(*bool) = r.binding.Enabled
	*dest[6].(**uuid.UUID) = r.binding.BotUserID
	*dest[7].(*json.RawMessage) = append(json.RawMessage(nil), r.binding.BotConfig...)
	*dest[8].(*time.Time) = r.binding.CreatedAt
	*dest[9].(*time.Time) = r.binding.UpdatedAt

	return nil
}

func (r *patchRows) Values() ([]any, error) {
	return nil, nil
}

func (r *patchRows) RawValues() [][]byte {
	return nil
}

func (r *patchRows) Conn() *pgx.Conn {
	return nil
}

func TestPatchAppliesEnabledAndMergesObjectPatch(t *testing.T) {
	bindingID := uuid.New()
	executor := &patchExecutor{
		queryFn: func(_ context.Context, query string, args ...any) (pgx.Rows, error) {
			if query != patchQuery {
				t.Fatalf("query = %q, want patchQuery", query)
			}
			if len(args) != 3 {
				t.Fatalf("query args = %d, want 3", len(args))
			}
			if args[0] != bindingID {
				t.Fatalf("binding ID argument = %v, want %s", args[0], bindingID)
			}
			enabled, ok := args[1].(bool)
			if !ok || enabled {
				t.Fatalf("enabled argument = %#v, want false", args[1])
			}
			patch, ok := args[2].(json.RawMessage)
			if !ok {
				t.Fatalf("bot config patch argument = %T, want json.RawMessage", args[2])
			}

			var existing map[string]json.RawMessage
			if err := json.Unmarshal([]byte(`{"bot_id":"twir-bot","keep":"value"}`), &existing); err != nil {
				t.Fatalf("unmarshal existing config: %v", err)
			}
			var update map[string]json.RawMessage
			if err := json.Unmarshal(patch, &update); err != nil {
				t.Fatalf("unmarshal patch: %v", err)
			}
			for key, value := range update {
				existing[key] = value
			}
			mergedConfig, err := json.Marshal(existing)
			if err != nil {
				t.Fatalf("marshal merged config: %v", err)
			}

			return &patchRows{binding: &model.ChannelPlatform{
				ID:                bindingID,
				ChannelID:         uuid.New(),
				Platform:          platform.PlatformTwitch,
				UserID:            uuid.New(),
				PlatformChannelID: "twitch-channel",
				Enabled:           enabled,
				BotConfig:         mergedConfig,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}}, nil
		},
	}
	repository := &Pgx{pool: executor}
	disabled := false

	binding, err := repository.Patch(context.Background(), bindingID, channelplatforms.PatchInput{
		Enabled:        &disabled,
		BotConfigPatch: json.RawMessage(`{"is_bot_mod":true}`),
	})
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}
	if binding.Enabled {
		t.Fatal("patched binding enabled = true, want false")
	}
	if executor.queryCalls != 1 {
		t.Fatalf("query calls = %d, want 1", executor.queryCalls)
	}

	var botConfig map[string]any
	if err := json.Unmarshal(binding.BotConfig, &botConfig); err != nil {
		t.Fatalf("unmarshal patched config: %v", err)
	}
	if botConfig["bot_id"] != "twir-bot" || botConfig["keep"] != "value" || botConfig["is_bot_mod"] != true {
		t.Fatalf("patched bot config = %v, want retained keys and is_bot_mod=true", botConfig)
	}
}

func TestPatchRejectsNonObjectBotConfigPatchBeforeQuery(t *testing.T) {
	tests := []struct {
		name  string
		patch json.RawMessage
	}{
		{name: "null", patch: json.RawMessage(`null`)},
		{name: "boolean", patch: json.RawMessage(`true`)},
		{name: "string", patch: json.RawMessage(`"value"`)},
		{name: "array", patch: json.RawMessage(`[]`)},
		{name: "malformed", patch: json.RawMessage(`{"is_bot_mod":`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &patchExecutor{
				queryFn: func(context.Context, string, ...any) (pgx.Rows, error) {
					t.Fatal("Patch invoked SQL for an invalid bot config patch")
					return nil, nil
				},
			}
			repository := &Pgx{pool: executor}

			_, err := repository.Patch(context.Background(), uuid.New(), channelplatforms.PatchInput{
				BotConfigPatch: tt.patch,
			})
			if !errors.Is(err, channelplatforms.ErrInvalidBotConfigPatch) {
				t.Fatalf("Patch error = %v, want ErrInvalidBotConfigPatch", err)
			}
			if executor.queryCalls != 0 {
				t.Fatalf("query calls = %d, want 0", executor.queryCalls)
			}
		})
	}
}

func TestPatchMapsNoRowsToErrNotFound(t *testing.T) {
	executor := &patchExecutor{
		queryFn: func(context.Context, string, ...any) (pgx.Rows, error) {
			return &patchRows{}, nil
		},
	}
	repository := &Pgx{pool: executor}
	disabled := false

	_, err := repository.Patch(context.Background(), uuid.New(), channelplatforms.PatchInput{Enabled: &disabled})
	if !errors.Is(err, channelplatforms.ErrNotFound) {
		t.Fatalf("Patch error = %v, want ErrNotFound", err)
	}
	if executor.queryCalls != 1 {
		t.Fatalf("query calls = %d, want 1", executor.queryCalls)
	}
}

func TestPatchWrapsQueryError(t *testing.T) {
	databaseErr := errors.New("database unavailable")
	executor := &patchExecutor{
		queryFn: func(context.Context, string, ...any) (pgx.Rows, error) {
			return nil, databaseErr
		},
	}
	repository := &Pgx{pool: executor}
	disabled := false

	_, err := repository.Patch(context.Background(), uuid.New(), channelplatforms.PatchInput{Enabled: &disabled})
	if !errors.Is(err, databaseErr) {
		t.Fatalf("Patch error = %v, want wrapped database error", err)
	}
	if !strings.Contains(err.Error(), "patch channel platform binding") {
		t.Fatalf("Patch error = %q, want patch context", err)
	}
}

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
