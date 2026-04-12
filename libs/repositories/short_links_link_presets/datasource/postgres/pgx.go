package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/twirapp/twir/libs/repositories/short_links_link_presets"
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

func (c *Pgx) GetByLinkID(ctx context.Context, linkID string) ([]repo.LinkPreset, error) {
	query, args, err := sq.
		Select("id", "link_id", "preset_id", "created_at").
		From("short_links_link_presets").
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

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[repo.LinkPreset])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []repo.LinkPreset{}, nil
		}
		return nil, err
	}

	return items, nil
}

func (c *Pgx) GetByPresetID(ctx context.Context, presetID string) ([]repo.LinkPreset, error) {
	query, args, err := sq.
		Select("id", "link_id", "preset_id", "created_at").
		From("short_links_link_presets").
		Where(squirrel.Eq{"preset_id": presetID}).
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

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[repo.LinkPreset])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []repo.LinkPreset{}, nil
		}
		return nil, err
	}

	return items, nil
}

func (c *Pgx) GetLinksByPresetID(ctx context.Context, presetID string) ([]string, error) {
	query, args, err := sq.
		Select("link_id").
		From("short_links_link_presets").
		Where(squirrel.Eq{"preset_id": presetID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var linkIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		linkIDs = append(linkIDs, id)
	}

	return linkIDs, nil
}

func (c *Pgx) Create(ctx context.Context, input repo.CreateInput) (repo.LinkPreset, error) {
	query, args, err := sq.
		Insert("short_links_link_presets").
		Columns("link_id", "preset_id").
		Values(input.LinkID, input.PresetID).
		Suffix("RETURNING id, link_id, preset_id, created_at").
		ToSql()
	if err != nil {
		return repo.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return repo.Nil, mapError(err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repo.LinkPreset])
	if err != nil {
		return repo.Nil, mapError(err)
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id string) error {
	query, args, err := sq.
		Delete("short_links_link_presets").
		Where(squirrel.Eq{"id": id}).
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

func (c *Pgx) DeleteByLinkAndPreset(ctx context.Context, linkID string, presetID string) error {
	query, args, err := sq.
		Delete("short_links_link_presets").
		Where(squirrel.Eq{"link_id": linkID, "preset_id": presetID}).
		ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	return err
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
