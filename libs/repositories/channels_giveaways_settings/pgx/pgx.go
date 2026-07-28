package pgx

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/channels_giveaways_settings"
	channels_giveaways_settings_repo "github.com/twirapp/twir/libs/repositories/channels_giveaways_settings"
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

var _ channels_giveaways_settings_repo.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) (channels_giveaways_settings.Settings, error) {
	query := `
INSERT INTO channels_giveaways_settings (channel_id)
VALUES ($1)
ON CONFLICT (channel_id) DO UPDATE SET channel_id = EXCLUDED.channel_id
RETURNING id, channel_id, created_at, updated_at, winner_message
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return channels_giveaways_settings.Nil, fmt.Errorf(
			"failed to get channels_giveaways_settings by channel_id: %w",
			err,
		)
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[channels_giveaways_settings.Settings],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return channels_giveaways_settings.Nil, nil
		}

		return channels_giveaways_settings.Nil, fmt.Errorf(
			"failed to collect channels_giveaways_settings by channel_id: %w",
			err,
		)
	}

	return result, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	channelID string,
	settings channels_giveaways_settings.Settings,
) (channels_giveaways_settings.Settings, error) {
	query := `
UPDATE channels_giveaways_settings
SET winner_message = $2, updated_at = NOW()
WHERE channel_id = $1
RETURNING id, channel_id, created_at, updated_at, winner_message
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID, settings.WinnerMessage)
	if err != nil {
		return channels_giveaways_settings.Nil, fmt.Errorf(
			"failed to update channels_giveaways_settings: %w",
			err,
		)
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[channels_giveaways_settings.Settings],
	)
	if err != nil {
		return channels_giveaways_settings.Nil, fmt.Errorf(
			"failed to collect updated channels_giveaways_settings: %w",
			err,
		)
	}

	return result, nil
}
