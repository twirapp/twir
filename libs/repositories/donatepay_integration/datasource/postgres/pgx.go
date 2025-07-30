package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	donatepayintegration "github.com/twirapp/twir/libs/repositories/donatepay_integration"
	"github.com/twirapp/twir/libs/repositories/donatepay_integration/model"
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

var _ donatepayintegration.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	model.DonatePayIntegration,
	error,
) {
	query := `
SELECT id, channel_id, api_key from channels_integrations_donatepay
WHERE channel_id = $1
LIMIT 1;
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return model.DonatePayIntegration{}, err
	}

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.DonatePayIntegration])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DonatePayIntegration{}, donatepayintegration.ErrNotFound
		}
		return model.DonatePayIntegration{}, err
	}

	return data, nil
}

func (c *Pgx) CreateOrUpdate(ctx context.Context, channelID, apiKey string) (
	model.DonatePayIntegration,
	error,
) {
	query := `
INSERT INTO channels_integrations_donatepay (channel_id, api_key)
VALUES ($1, $2)
ON CONFLICT (channel_id) DO UPDATE
SET api_key = $2
RETURNING id, channel_id, api_key;
`

	rows, err := c.pool.Query(ctx, query, channelID, apiKey)
	if err != nil {
		return model.DonatePayIntegration{}, err
	}

	data, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.DonatePayIntegration])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DonatePayIntegration{}, donatepayintegration.ErrNotFound
		}
		return model.DonatePayIntegration{}, err
	}

	return data, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_integrations_donatepay
WHERE id = $1;
`

	_, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return donatepayintegration.ErrNotFound
		}
		return err
	}

	return nil
}
