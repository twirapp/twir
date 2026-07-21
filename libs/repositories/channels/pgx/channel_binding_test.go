package pgx

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCollectExactlyOneChannelRejectsMultipleRows(t *testing.T) {
	_, err := collectExactlyOneChannel(&channelRowsFixture{rows: []channelRowFixture{
		{id: uuid.New(), bindings: []byte(`[]`)},
		{id: uuid.New(), bindings: []byte(`[]`)},
	}})
	if !errors.Is(err, pgx.ErrTooManyRows) {
		t.Fatalf("collect exactly one channel error = %v, want ErrTooManyRows", err)
	}
}

func TestCreateChannelAndBindingsQueryIsAtomic(t *testing.T) {
	assertAtomicChannelBindingQuery(t, createChannelAndBindingsQuery, "INSERT INTO channels", "created_channel")

	for _, fragment := range []string{
		"twitch_bot_enabled = false",
		"kick_bot_enabled = false",
		"kick_bot_id IS NULL",
	} {
		if !strings.Contains(createChannelAndBindingsQuery, fragment) {
			t.Fatalf("create query does not reject orphaned provider state %q", fragment)
		}
	}
}

func TestUpdateChannelAndBindingsQueryIsAtomic(t *testing.T) {
	assertAtomicChannelBindingQuery(t, updateChannelAndBindingsQuery, "UPDATE channels", "updated_channel")

	if got := strings.Count(updateChannelAndBindingsQuery, "ON CONFLICT (channel_id, platform) DO UPDATE"); got != 2 {
		t.Fatalf("update binding upserts = %d, want 2", got)
	}
}

func TestGetAllByBindingPlatformQueryIsComplete(t *testing.T) {
	for _, fragment := range []string{
		"WHERE EXISTS",
		"FROM channel_platforms cp_filter",
		"cp_filter.channel_id = c.id",
		"cp_filter.platform = $1",
		"ORDER BY c.id",
	} {
		if !strings.Contains(getAllByBindingPlatformQuery, fragment) {
			t.Fatalf("platform binding query does not contain %q", fragment)
		}
	}

	if strings.Contains(strings.ToUpper(getAllByBindingPlatformQuery), "LIMIT") {
		t.Fatal("platform binding query must not apply a result limit")
	}
}

func assertAtomicChannelBindingQuery(t *testing.T, query, channelMutation, channelCTE string) {
	t.Helper()

	for _, fragment := range []string{
		"WITH",
		channelMutation,
		"INSERT INTO channel_platforms",
		"SELECT id FROM " + channelCTE,
	} {
		if !strings.Contains(query, fragment) {
			t.Fatalf("query does not contain %q", fragment)
		}
	}

	if got := strings.Count(query, "INSERT INTO channel_platforms"); got != 2 {
		t.Fatalf("binding inserts = %d, want 2", got)
	}
}

type channelRowFixture struct {
	id       uuid.UUID
	apiKey   *string
	bindings []byte
}

type channelRowsFixture struct {
	rows  []channelRowFixture
	index int
}

func (r *channelRowsFixture) Close() {}

func (r *channelRowsFixture) Err() error {
	return nil
}

func (r *channelRowsFixture) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("")
}

func (r *channelRowsFixture) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r *channelRowsFixture) Next() bool {
	if r.index == len(r.rows) {
		return false
	}

	r.index++
	return true
}

func (r *channelRowsFixture) Scan(dest ...any) error {
	if len(dest) != 3 {
		return fmt.Errorf("scan destinations = %d, want 3", len(dest))
	}

	row := r.rows[r.index-1]
	channelID, ok := dest[0].(*uuid.UUID)
	if !ok {
		return fmt.Errorf("channel ID destination = %T, want *uuid.UUID", dest[0])
	}
	apiKey, ok := dest[1].(**string)
	if !ok {
		return fmt.Errorf("API key destination = %T, want **string", dest[1])
	}
	bindings, ok := dest[2].(*[]byte)
	if !ok {
		return fmt.Errorf("bindings destination = %T, want *[]byte", dest[2])
	}

	*channelID = row.id
	*apiKey = row.apiKey
	*bindings = append((*bindings)[:0], row.bindings...)
	return nil
}

func (r *channelRowsFixture) Values() ([]any, error) {
	return nil, nil
}

func (r *channelRowsFixture) RawValues() [][]byte {
	return nil
}

func (r *channelRowsFixture) Conn() *pgx.Conn {
	return nil
}
