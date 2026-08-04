package postgres

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	provider "github.com/twirapp/twir/libs/integrations/streamlabs"
	repository "github.com/twirapp/twir/libs/repositories/streamlabs_integration"
)

func TestTokenStoreReadsEnabledCredentials(t *testing.T) {
	t.Parallel()

	db := &fakeTokenStoreDB{row: fakeRow{values: []string{"access", "refresh"}}}
	repo := &Pgx{tokenStoreDB: db}

	tokens, err := repo.GetTokens(context.Background(), "channel")
	if err != nil {
		t.Fatalf("GetTokens() error = %v", err)
	}
	if want := (provider.Tokens{AccessToken: "access", RefreshToken: "refresh"}); tokens != want {
		t.Fatalf("GetTokens() = %#v, want %#v", tokens, want)
	}
	if !strings.Contains(db.query, "enabled = TRUE") {
		t.Fatalf("GetTokens() query does not reject disabled rows: %s", db.query)
	}
	if !reflect.DeepEqual(db.args, []any{"channel"}) {
		t.Fatalf("GetTokens() args = %#v, want channel", db.args)
	}
}

func TestTokenStoreUpdatesBothCredentialsInOneStatement(t *testing.T) {
	t.Parallel()

	db := &fakeTokenStoreDB{tag: pgconn.NewCommandTag("UPDATE 1")}
	repo := &Pgx{tokenStoreDB: db}
	tokens := provider.Tokens{AccessToken: "new-access", RefreshToken: "new-refresh"}

	if err := repo.UpdateTokens(context.Background(), "channel", tokens); err != nil {
		t.Fatalf("UpdateTokens() error = %v", err)
	}
	if db.execCalls != 1 {
		t.Fatalf("UpdateTokens() statements = %d, want 1", db.execCalls)
	}
	if !strings.Contains(db.query, `"access_token" = $2`) ||
		!strings.Contains(db.query, `"refresh_token" = $3`) ||
		!strings.Contains(db.query, "enabled = TRUE") {
		t.Fatalf("UpdateTokens() query does not atomically update enabled credentials: %s", db.query)
	}
	if want := []any{"channel", "new-access", "new-refresh"}; !reflect.DeepEqual(db.args, want) {
		t.Fatalf("UpdateTokens() args = %#v, want %#v", db.args, want)
	}
}

func TestTokenStoreReadRejectsDisabledRow(t *testing.T) {
	t.Parallel()

	db := &fakeTokenStoreDB{row: fakeRow{err: pgx.ErrNoRows}}
	repo := &Pgx{tokenStoreDB: db}

	_, err := repo.GetTokens(context.Background(), "channel")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("GetTokens() error = %v, want ErrNotFound for disabled row", err)
	}
}

func TestTokenStoreUpdateRejectsDisabledRow(t *testing.T) {
	t.Parallel()

	db := &fakeTokenStoreDB{tag: pgconn.NewCommandTag("UPDATE 0")}
	repo := &Pgx{tokenStoreDB: db}
	err := repo.UpdateTokens(context.Background(), "channel", provider.Tokens{
		AccessToken: "new-access", RefreshToken: "new-refresh",
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("UpdateTokens() error = %v, want ErrNotFound for disabled row", err)
	}
}

func TestTokenStoreRefusesIncompleteCredentialsWithoutWriting(t *testing.T) {
	t.Parallel()

	db := &fakeTokenStoreDB{tag: pgconn.NewCommandTag("UPDATE 1")}
	repo := &Pgx{tokenStoreDB: db}
	err := repo.UpdateTokens(context.Background(), "channel", provider.Tokens{
		AccessToken: "new-access",
	})
	if err == nil {
		t.Fatal("UpdateTokens() error = nil, want incomplete-credential rejection")
	}
	if db.execCalls != 0 {
		t.Fatalf("UpdateTokens() statements = %d, want none", db.execCalls)
	}
}

type fakeTokenStoreDB struct {
	row       pgx.Row
	tag       pgconn.CommandTag
	query     string
	args      []any
	execCalls int
}

func (f *fakeTokenStoreDB) QueryRow(_ context.Context, query string, args ...any) pgx.Row {
	f.query = query
	f.args = args
	return f.row
}

func (f *fakeTokenStoreDB) Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	f.execCalls++
	f.query = query
	f.args = args
	return f.tag, nil
}

type fakeRow struct {
	values []string
	err    error
}

func (f fakeRow) Scan(dest ...any) error {
	if f.err != nil {
		return f.err
	}
	for index, value := range f.values {
		pointer, ok := dest[index].(*string)
		if !ok {
			return errors.New("unexpected scan destination")
		}
		*pointer = value
	}
	return nil
}
