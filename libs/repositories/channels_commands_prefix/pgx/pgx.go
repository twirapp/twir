package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
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

var _ channels_commands_prefix.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	model.ChannelsCommandsPrefix,
	error,
) {
	query := `
SELECT id, channel_id, prefix, created_at, updated_at from channels_commands_prefix
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to get channels_commands_prefix by channel_id: %w", err)
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelsCommandsPrefix])
	if err != nil {
		return model.Nil, fmt.Errorf(
			"failed to collect channels_commands_prefix by channel_id: %w",
			err,
		)
	}

	return data, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input channels_commands_prefix.CreateInput,
) (model.ChannelsCommandsPrefix, error) {
	query := `
INSERT INTO channels_commands_prefix (channel_id, prefix)
VALUES ($1, $2)
RETURNING id, channel_id, prefix, created_at, updated_at
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.ChannelID, input.Prefix)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to create channels_commands_prefix: %w", err)
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelsCommandsPrefix])
	if err != nil {
		return model.Nil, fmt.Errorf("failed to collect created channels_commands_prefix: %w", err)
	}

	return data, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input channels_commands_prefix.UpdateInput,
) (model.ChannelsCommandsPrefix, error) {
	query := `
UPDATE channels_commands_prefix
SET prefix = $1, updated_at = now()
WHERE id = $2
RETURNING id, channel_id, prefix, created_at, updated_at
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.Prefix, id)
	if err != nil {
		return model.Nil, fmt.Errorf("failed to update channels_commands_prefix: %w", err)
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelsCommandsPrefix])
	if err != nil {
		return model.Nil, fmt.Errorf("failed to collect updated channels_commands_prefix: %w", err)
	}

	return data, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_commands_prefix
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete channels_commands_prefix: %w", err)
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("failed to delete channels_commands_prefix: no rows affected")
	}

	return nil
}
