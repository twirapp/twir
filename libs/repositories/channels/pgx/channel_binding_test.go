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

func TestCreateChannelQueryDoesNotWriteProviderLinkage(t *testing.T) {
	for _, legacyColumn := range []string{
		"twitch_user_id",
		"kick_user_id",
		"twitch_bot_enabled",
		"kick_bot_enabled",
		"kick_bot_id",
		"botId",
		"channel_platforms",
	} {
		if strings.Contains(createChannelQuery, legacyColumn) {
			t.Fatalf("generic channel creation writes legacy provider linkage %q", legacyColumn)
		}
	}

	if !strings.Contains(createChannelQuery, `INSERT INTO channels DEFAULT VALUES`) {
		t.Fatalf("generic channel creation query = %q, want independent channel insert", createChannelQuery)
	}
}

func TestUpdateChannelQueryUpdatesMutableChannelSettings(t *testing.T) {
	for _, fragment := range []string{
		"UPDATE channels",
		`SET "isEnabled" = COALESCE($2, "isEnabled")`,
		"api_key = COALESCE($3, api_key)",
		"WHERE id = $1",
		"RETURNING id",
	} {
		if !strings.Contains(updateChannelQuery, fragment) {
			t.Fatalf("update channel query does not contain %q", fragment)
		}
	}

	for _, legacyColumn := range []string{
		"twitch_user_id",
		"kick_user_id",
		"twitch_bot_enabled",
		"kick_bot_enabled",
		"kick_bot_id",
		"botId",
		"channel_platforms",
	} {
		if strings.Contains(updateChannelQuery, legacyColumn) {
			t.Fatalf("update channel query writes legacy provider linkage %q", legacyColumn)
		}
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

type channelRowFixture struct {
	id       uuid.UUID
	apiKey   *string
	planID   *string
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
	if len(dest) != 4 {
		return fmt.Errorf("scan destinations = %d, want 4", len(dest))
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
	planID, ok := dest[2].(**string)
	if !ok {
		return fmt.Errorf("plan ID destination = %T, want **string", dest[2])
	}
	bindings, ok := dest[3].(*[]byte)
	if !ok {
		return fmt.Errorf("bindings destination = %T, want *[]byte", dest[3])
	}

	*channelID = row.id
	*apiKey = row.apiKey
	*planID = row.planID
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
