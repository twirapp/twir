package pgx

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestAssignVKVideoLiveBotUpdatesOnlyVKBindingsWithoutTheBot(t *testing.T) {
	// Given
	botUserID := uuid.New()
	affectedChannelIDs := []uuid.UUID{uuid.New(), uuid.New()}
	executor := &vkVideoBotAssignmentExecutor{
		queryFn: func(_ context.Context, query string, args ...any) (pgx.Rows, error) {
			lowerQuery := strings.Join(strings.Fields(strings.ToLower(query)), " ")
			for _, fragment := range []string{
				"update channel_platforms",
				"set bot_user_id = $1",
				"where platform = $2",
				"bot_user_id is distinct from $1",
				"returning channel_id",
			} {
				if !strings.Contains(lowerQuery, fragment) {
					t.Fatalf("assignment query missing %q: %s", fragment, query)
				}
			}
			for _, forbidden := range []string{
				"set user_id =",
				", user_id =",
				"set platform_channel_id =",
				", platform_channel_id =",
				"set enabled =",
				", enabled =",
				"set bot_config =",
				", bot_config =",
			} {
				if strings.Contains(lowerQuery, forbidden) {
					t.Fatalf("assignment query overwrites unrelated binding field %q: %s", forbidden, query)
				}
			}
			if len(args) != 2 || args[0] != botUserID || args[1] != platform.PlatformVKVideoLive {
				t.Fatalf("assignment args = %#v, want bot user ID and VK Video Live platform", args)
			}

			return &vkVideoBotAssignmentRows{channelIDs: affectedChannelIDs}, nil
		},
	}
	repository := &Pgx{pool: executor}

	// When
	updatedChannelIDs, err := repository.AssignVKVideoLiveBot(context.Background(), botUserID)

	// Then
	if err != nil {
		t.Fatalf("assign VK Video Live bot: %v", err)
	}
	if len(updatedChannelIDs) != len(affectedChannelIDs) {
		t.Fatalf("updated channel ID count = %d, want %d", len(updatedChannelIDs), len(affectedChannelIDs))
	}
	for index, channelID := range affectedChannelIDs {
		if updatedChannelIDs[index] != channelID {
			t.Fatalf("updatedChannelIDs[%d] = %s, want %s", index, updatedChannelIDs[index], channelID)
		}
	}
}

type vkVideoBotAssignmentExecutor struct {
	queryFn func(context.Context, string, ...any) (pgx.Rows, error)
}

func (e *vkVideoBotAssignmentExecutor) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return e.queryFn(ctx, query, args...)
}

func (*vkVideoBotAssignmentExecutor) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}

type vkVideoBotAssignmentRows struct {
	channelIDs []uuid.UUID
	index      int
}

func (*vkVideoBotAssignmentRows) Close() {}

func (*vkVideoBotAssignmentRows) Err() error { return nil }

func (*vkVideoBotAssignmentRows) CommandTag() pgconn.CommandTag { return pgconn.NewCommandTag("") }

func (*vkVideoBotAssignmentRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{{Name: "channel_id"}}
}

func (r *vkVideoBotAssignmentRows) Next() bool {
	if r.index >= len(r.channelIDs) {
		return false
	}
	r.index++
	return true
}

func (r *vkVideoBotAssignmentRows) Scan(dest ...any) error {
	if len(dest) != 1 {
		return fmt.Errorf("scan destinations = %d, want 1", len(dest))
	}
	*dest[0].(*uuid.UUID) = r.channelIDs[r.index-1]
	return nil
}

func (*vkVideoBotAssignmentRows) Values() ([]any, error) { return nil, nil }

func (*vkVideoBotAssignmentRows) RawValues() [][]byte { return nil }

func (*vkVideoBotAssignmentRows) Conn() *pgx.Conn { return nil }
