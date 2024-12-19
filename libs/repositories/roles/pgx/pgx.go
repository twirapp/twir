package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/roles"
	"github.com/twirapp/twir/libs/repositories/roles/model"
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

var _ roles.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetManyByIDS(ctx context.Context, ids []uuid.UUID) ([]model.Role, error) {
	query := `
SELECT id, "channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time
FROM channels_roles
WHERE id = ANY($1)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("GetManyByIDS: failed to execute select query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return nil, fmt.Errorf("GetManyByIDS: failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Role, error) {
	query := `
SELECT id, "channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time
FROM channels_roles
WHERE "channelId" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("GetManyByChannelID: failed to execute select query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return nil, fmt.Errorf("GetManyByChannelID: failed to collect rows: %w", err)
	}

	return result, nil
}
