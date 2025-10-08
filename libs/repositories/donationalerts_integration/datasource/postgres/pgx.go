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
SELECT ci.id, ci.enabled, ci."channelId", ci."integrationId", ci."accessToken",
       ci."refreshToken", ci."clientId", ci."clientSecret", ci."apiKey",
       ci.data
FROM channels_integrations ci
JOIN integrations i ON ci."integrationId" = i.id
WHERE ci."channelId" = $1 AND i.service = 'DONATIONALERTS'
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
UPDATE channels_integrations as ci
SET
	"enabled" = COALESCE($2, ci.enabled),
	"accessToken" = COALESCE($3, ci."accessToken"),
	"refreshToken" = COALESCE($4, ci."refreshToken"),
	"clientId" = COALESCE($5, ci."clientId"),
	"clientSecret" = COALESCE($6, ci."clientSecret"),
	"apiKey" = COALESCE($7, ci."apiKey"),
	data = COALESCE($8, ci.data)
FROM integrations i
WHERE ci."channelId" = $1
  AND ci."integrationId" = i.id
  AND i.service = 'DONATIONALERTS';
`

	_, err := c.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.Enabled,
		opts.AccessToken,
		opts.RefreshToken,
		opts.ClientID,
		opts.ClientSecret,
		opts.APIKey,
		opts.Data,
	)

	return err
}

func (c *Pgx) Delete(ctx context.Context, channelID string) error {
	query := `
DELETE FROM channels_integrations
USING integrations i
WHERE channels_integrations."channelId" = $1
  AND channels_integrations."integrationId" = i.id
  AND i.service = 'DONATIONALERTS';
`

	_, err := c.pool.Exec(ctx, query, channelID)
	return err
}

func (c *Pgx) Create(
	ctx context.Context,
	opts donationalertsintegration.CreateOpts,
) error {
	query := `
INSERT INTO channels_integrations ("channelId", "integrationId", enabled, "accessToken",
                                  "refreshToken", "clientId", "clientSecret", "apiKey", data)
SELECT $1, i.id, $2, $3, $4, $5, $6, $7, $8
FROM integrations i
WHERE i.service = 'DONATIONALERTS';
`

	_, err := c.pool.Exec(
		ctx,
		query,
		opts.ChannelID,
		opts.Enabled,
		opts.AccessToken,
		opts.RefreshToken,
		opts.ClientID,
		opts.ClientSecret,
		opts.APIKey,
		opts.Data,
	)

	return err
}
