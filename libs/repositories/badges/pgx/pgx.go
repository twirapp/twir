package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/badges"
	"github.com/twirapp/twir/libs/repositories/badges/model"
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

var _ badges.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetMany(ctx context.Context, input badges.GetManyInput) ([]model.Badge, error) {
	query := `
SELECT id, name, enabled, created_at, file_name, ffz_slot
FROM badges
WHERE enabled = $1
`

	rows, err := c.pool.Query(ctx, query, input.Enabled)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Badge])
	if err != nil {
		return nil, err
	}

	return result, nil
}
