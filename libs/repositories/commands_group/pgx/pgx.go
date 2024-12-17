package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/commands_group"
	"github.com/twirapp/twir/libs/repositories/commands_group/model"
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

var _ commands_group.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetManyByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Group, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := `
SELECT id, "channelId", name, color
FROM channels_commands_groups
WHERE id IN ($1)
`

	idsStrings := make([]string, len(ids))
	for i, id := range ids {
		idsStrings[i] = id.String()
	}

	rows, err := c.pool.Query(ctx, query, idsStrings)
	if err != nil {
		return nil, err
	}

	groups, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Group])
	if err != nil {
		return nil, err
	}

	return groups, nil
}
