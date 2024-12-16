package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/badges-users"
	"github.com/twirapp/twir/libs/repositories/badges-users/model"
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

var _ badges_users.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetMany(ctx context.Context, input badges_users.GetManyInput) (
	[]model.BadgeUser,
	error,
) {
	query := `
SELECT id, badge_id, user_id, created_at
FROM badges_users
WHERE badge_id = $1
`

	rows, err := c.pool.Query(ctx, query, input.BadgeID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.BadgeUser])
	if err != nil {
		return nil, err
	}

	return result, nil
}
