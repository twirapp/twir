package postgres

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	donationalertsintegration "github.com/twirapp/twir/libs/repositories/donationalerts_integration"
	"github.com/twirapp/twir/libs/repositories/donationalerts_integration/model"
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

var _ donationalertsintegration.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	model.DonationAlertsIntegration,
	error,
) {
	query := `
SELECT id, public_id, channel_id, access_token, refresh_token, username, avatar, created_at, updated_at, enabled
FROM channels_integrations_donationalerts
WHERE channel_id = $1
LIMIT 1;
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, err
	}

	data, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.DonationAlertsIntegration],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, donationalertsintegration.ErrNotFound
		}
		return model.Nil, err
	}

	return data, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	opts donationalertsintegration.UpdateOpts,
) error {
	query := `
UPDATE channels_integrations_donationalerts
SET
	"enabled" = COALESCE($2, enabled),
	"access_token" = COALESCE($3, "access_token"),
	"refresh_token" = COALESCE($4, "refresh_token"),
	username = COALESCE($5, username),
	avatar = COALESCE($6, avatar),
	updated_at = NOW()
WHERE channel_id = $1
`

	cmd, err := c.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.Enabled,
		opts.AccessToken,
		opts.RefreshToken,
		opts.UserName,
		opts.Avatar,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return donationalertsintegration.ErrNotFound
	}

	return err
}

func (c *Pgx) Delete(ctx context.Context, channelID string) error {
	query := `
DELETE FROM channels_integrations_donationalerts
WHERE channel_id = $1
`

	_, err := c.pool.Exec(ctx, query, channelID)
	return err
}

func (c *Pgx) Create(
	ctx context.Context,
	opts donationalertsintegration.CreateOpts,
) error {
	query := `
INSERT INTO channels_integrations_donationalerts (channel_id, enabled, access_token, refresh_token, username, avatar)
VALUES ($1, $2, $3, $4, $5, $6)
`

	_, err := c.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.Enabled,
		opts.AccessToken,
		opts.RefreshToken,
		opts.UserName,
		opts.Avatar,
	)

	return err
}
