package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	twitchconduits "github.com/twirapp/twir/libs/repositories/twitch_conduits"
	"github.com/twirapp/twir/libs/repositories/twitch_conduits/model"
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

var _ twitchconduits.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetOne(ctx context.Context) (model.Conduit, error) {
	query := `
SELECT id, created_at, updated_at, shard_count
FROM twitch_conduits
LIMIT 1;
`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return model.Nil, err
	}

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Conduit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}

		return model.Nil, err
	}

	return data, nil
}

func (c *Pgx) Create(ctx context.Context, input twitchconduits.CreateInput) (model.Conduit, error) {
	query := `
INSERT INTO twitch_conduits (id, shard_count)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, shard_count;
`

	rows, err := c.pool.Query(ctx, query, input.ID, input.ShardCount)
	if err != nil {
		return model.Nil, err
	}

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Conduit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}

		return model.Nil, err
	}

	return data, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id string,
	input twitchconduits.UpdateInput,
) (model.Conduit, error) {
	query := `
UPDATE twitch_conduits
SET shard_count = $1, updated_at = NOW()
WHERE id = $2
RETURNING id, created_at, updated_at, shard_count;
`

	rows, err := c.pool.Query(ctx, query, input.ShardCount, id)
	if err != nil {
		return model.Nil, err
	}

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Conduit])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}

		return model.Nil, err
	}

	return data, nil
}

func (c *Pgx) Delete(ctx context.Context, id string) error {
	query := `
DELETE FROM twitch_conduits
WHERE id = $1;
`

	_, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	return nil
}

func (c *Pgx) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM twitch_conduits;`
	_, err := c.pool.Exec(ctx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}
