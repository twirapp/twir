package pgx

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestLockByChannelIDConsumesEveryBindingInDeterministicOrder(t *testing.T) {
	channelID := uuid.New()
	rows := &lockRows{ids: []uuid.UUID{uuid.New(), uuid.New()}}
	executor := &lockExecutor{queryFn: func(_ context.Context, query string, args ...any) (pgx.Rows, error) {
		if len(args) != 1 || args[0] != channelID {
			t.Fatalf("lock query args = %#v, want channel ID %s", args, channelID)
		}
		for _, fragment := range []string{"select id", "from channel_platforms", "where channel_id = $1", "order by id", "for update"} {
			if !strings.Contains(strings.ToLower(query), fragment) {
				t.Fatalf("lock query missing %q: %s", fragment, query)
			}
		}
		return rows, nil
	}}
	repository := &Pgx{pool: executor}

	if err := repository.LockByChannelID(context.Background(), channelID); err != nil {
		t.Fatalf("LockByChannelID() error = %v", err)
	}
	if rows.nextCalls != len(rows.ids)+1 {
		t.Fatalf("lock rows Next calls = %d, want %d to consume all rows", rows.nextCalls, len(rows.ids)+1)
	}
	if !rows.closed {
		t.Fatal("LockByChannelID() did not close fully consumed rows")
	}
}

type lockExecutor struct {
	queryFn func(context.Context, string, ...any) (pgx.Rows, error)
}

func (e *lockExecutor) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return e.queryFn(ctx, query, args...)
}

func (*lockExecutor) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}

type lockRows struct {
	ids       []uuid.UUID
	index     int
	nextCalls int
	closed    bool
}

func (r *lockRows) Close() { r.closed = true }

func (*lockRows) Err() error { return nil }

func (*lockRows) CommandTag() pgconn.CommandTag { return pgconn.NewCommandTag("") }

func (*lockRows) FieldDescriptions() []pgconn.FieldDescription {
	return []pgconn.FieldDescription{{Name: "id"}}
}

func (r *lockRows) Next() bool {
	r.nextCalls++
	if r.index >= len(r.ids) {
		return false
	}
	r.index++
	return true
}

func (r *lockRows) Scan(dest ...any) error {
	if len(dest) != 1 {
		return fmt.Errorf("scan destinations = %d, want 1", len(dest))
	}
	*dest[0].(*uuid.UUID) = r.ids[r.index-1]
	return nil
}

func (*lockRows) Values() ([]any, error) { return nil, nil }

func (*lockRows) RawValues() [][]byte { return nil }

func (*lockRows) Conn() *pgx.Conn { return nil }
