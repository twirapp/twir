package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/users"
	"github.com/twirapp/twir/libs/repositories/users/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ users.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetByID(ctx context.Context, id string) (model.User, error) {
	query := `
SELECT id, "tokenId", "isTester", "isBotAdmin", "apiKey", is_banned, hide_on_landing_page
FROM users
WHERE id = $1
`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return model.Nil, err
	}

	return user, nil
}

func (c *Pgx) GetManyByIDS(ctx context.Context, input users.GetManyInput) ([]model.User, error) {
	selectBuilder := sq.Select(
		"id",
		`"tokenId"`,
		`"isTester"`,
		"isBotAdmin",
		"apiKey",
		"is_banned",
		"hide_on_landing_page",
	).From("users")

	limit := 100
	if input.Limit > 0 {
		limit = input.Limit
	}

	selectBuilder = selectBuilder.Limit(uint64(limit))

	perPage := 100
	if input.PerPage > 0 {
		perPage = input.PerPage
	}

	selectBuilder = selectBuilder.Offset(uint64(perPage))

	if len(input.IDs) > 0 {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"id": input.IDs})
	}

	if input.IsBotAdmin != nil {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"isBotAdmin": input.IsBotAdmin})
	}

	if input.IsBanned != nil {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"is_banned": input.IsBanned})
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id string, input users.UpdateInput) (model.User, error) {
	updateBuilder := sq.Update("users").Where(squirrel.Eq{"id": id})

	if input.IsBotAdmin != nil {
		updateBuilder = updateBuilder.Set(`"isBotAdmin"`, input.IsBotAdmin)
	}

	if input.IsBanned != nil {
		updateBuilder = updateBuilder.Set("is_banned", input.IsBanned)
	}

	if input.ApiKey != nil {
		updateBuilder = updateBuilder.Set(`"apiKey"`, input.ApiKey)
	}

	if input.HideOnLandingPage != nil {
		updateBuilder = updateBuilder.Set("hide_on_landing_page", input.HideOnLandingPage)
	}

	if input.TokenID != nil {
		updateBuilder = updateBuilder.Set(`"tokenId"`, input.TokenID)
	}

	updateBuilder = updateBuilder.Suffix(`RETURNING id, "tokenId", "isTester", "isBotAdmin", "apiKey", is_banned, hide_on_landing_page`)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return model.Nil, err
	}

	return user, nil
}
