package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/pastebins"
	"github.com/twirapp/twir/libs/repositories/pastebins/model"
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
	_  pastebins.Repository = (*Pgx)(nil)
	sq                      = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) Count(ctx context.Context, input pastebins.CountInput) (int64, error) {
	selectBuilder := sq.Select("COUNT(*)").From("pastebins")

	if input.OwnerUserID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"owner_user_id": input.OwnerUserID})
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	err = conn.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) Create(ctx context.Context, input pastebins.CreateInput) (model.Pastebin, error) {
	query := `
INSERT INTO pastebins (id, content, "expire_at", "owner_user_id", "user_ip", "user_agent")
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, content, "expire_at", "owner_user_id", "user_ip", "user_agent"
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.ID, input.Content, input.ExpireAt, input.OwnerUserID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Pastebin])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByID(ctx context.Context, id string) (model.Pastebin, error) {
	query := `
SELECT id, created_at, content, "expire_at", "owner_user_id", "user_ip", "user_agent"
FROM pastebins
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Pastebin])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, pastebins.ErrNotFound
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id string) error {
	query := `
DELETE FROM pastebins
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}

func (c *Pgx) GetManyByOwner(ctx context.Context, input pastebins.GetManyInput) (
	pastebins.GetManyOutput,
	error,
) {
	builder := sq.Select(
		"id",
		"created_at",
		"content",
		"expire_at",
		"owner_user_id",
		"user_ip",
		"user_agent",
	).From("pastebins").
		Where(squirrel.Eq{"owner_user_id": input.OwnerUserID}).
		OrderBy("created_at DESC")

	perPage := input.PerPage
	if perPage == 0 {
		perPage = 20
	}

	offset := input.Page * perPage

	if input.Page > 0 && input.PerPage > 0 {
		builder = builder.Limit(uint64(perPage)).Offset(uint64(offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return pastebins.GetManyOutput{}, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return pastebins.GetManyOutput{}, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Pastebin])
	if err != nil {
		return pastebins.GetManyOutput{}, err
	}

	countQuery := `
SELECT COUNT(*) FROM pastebins
WHERE "owner_user_id" = $1
`

	var count int
	err = conn.QueryRow(ctx, countQuery, input.OwnerUserID).Scan(&count)
	if err != nil {
		return pastebins.GetManyOutput{}, err
	}

	return pastebins.GetManyOutput{
		Items: result,
		Total: count,
	}, nil
}
