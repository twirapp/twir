package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_presets"
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

func (c *Pgx) GetByUserID(ctx context.Context, userID string) ([]repo.Preset, error) {
	query, args, err := sq.
		Select("id", "user_id", "name", "description", "created_at", "updated_at").
		From("short_links_banned_ua_presets").
		Where(squirrel.Eq{"user_id": userID}).
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

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[repo.Preset])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []repo.Preset{}, nil
		}
		return nil, err
	}

	return items, nil
}

func (c *Pgx) GetByID(ctx context.Context, id string) (repo.Preset, error) {
	query, args, err := sq.
		Select("id", "user_id", "name", "description", "created_at", "updated_at").
		From("short_links_banned_ua_presets").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return repo.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return repo.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repo.Preset])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Nil, repo.ErrNotFound
		}
		return repo.Nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input repo.CreateInput) (repo.Preset, error) {
	query, args, err := sq.
		Insert("short_links_banned_ua_presets").
		Columns("user_id", "name", "description").
		Values(input.UserID, input.Name, input.Description).
		Suffix("RETURNING id, user_id, name, description, created_at, updated_at").
		ToSql()
	if err != nil {
		return repo.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return repo.Nil, mapError(err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repo.Preset])
	if err != nil {
		return repo.Nil, mapError(err)
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id string, input repo.UpdateInput) (repo.Preset, error) {
	updateBuilder := sq.Update("short_links_banned_ua_presets").
		Where(squirrel.Eq{"id": id}).
		Set("updated_at", time.Now())

	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", *input.Name)
	}
	if input.Description != nil {
		updateBuilder = updateBuilder.Set("description", *input.Description)
	}

	query, args, err := updateBuilder.
		Suffix("RETURNING id, user_id, name, description, created_at, updated_at").
		ToSql()
	if err != nil {
		return repo.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return repo.Nil, mapError(err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[repo.Preset])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Nil, repo.ErrNotFound
		}
		return repo.Nil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id string, userID string) error {
	query, args, err := sq.
		Delete("short_links_banned_ua_presets").
		Where(squirrel.Eq{"id": id, "user_id": userID}).
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
