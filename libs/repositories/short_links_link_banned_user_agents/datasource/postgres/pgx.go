package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/twirapp/twir/libs/repositories/short_links_link_banned_user_agents"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var (
	_  repo.Repository = (*Pgx)(nil)
	sq                 = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByLinkID(ctx context.Context, linkID string) ([]repo.BannedUserAgent, error) {
	query, args, err := sq.
		Select("id", "link_id", "pattern", "description", "created_at").
		From("short_links_link_banned_user_agents").
		Where(squirrel.Eq{"link_id": linkID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[repo.BannedUserAgent])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []repo.BannedUserAgent{}, nil
		}
		return nil, err
	}

	return items, nil
}

func (c *Pgx) Create(ctx context.Context, input repo.CreateInput) (repo.BannedUserAgent, error) {
	query, args, err := sq.
		Insert("short_links_link_banned_user_agents").
		Columns("link_id", "pattern", "description").
		Values(input.LinkID, input.Pattern, input.Description).
		Suffix("RETURNING id, link_id, pattern, description, created_at").
		ToSql()
	if err != nil {
		return repo.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return repo.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repo.BannedUserAgent])
	if err != nil {
		return repo.Nil, mapError(err)
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id string, linkID string) error {
	query, args, err := sq.
		Delete("short_links_link_banned_user_agents").
		Where(squirrel.Eq{"id": id, "link_id": linkID}).
		ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	tag, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return repo.ErrNotFound
	}

	return nil
}

func mapError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code == "23505" {
		return repo.ErrAlreadyExists
	}

	return err
}
