package postgres

import (
	"context"
	"errors"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	provider "github.com/twirapp/twir/libs/integrations/streamlabs"
	streamlabsintegration "github.com/twirapp/twir/libs/repositories/streamlabs_integration"
	"github.com/twirapp/twir/libs/repositories/streamlabs_integration/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:         opts.PgxPool,
		getter:       trmpgx.DefaultCtxGetter,
		tokenStoreDB: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ streamlabsintegration.Repository = (*Pgx)(nil)
var _ provider.TokenStore = (*Pgx)(nil)

type tokenStoreDB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type Pgx struct {
	pool         *pgxpool.Pool
	getter       *trmpgx.CtxGetter
	tokenStoreDB tokenStoreDB
}

func (c *Pgx) GetTokens(ctx context.Context, channelID string) (provider.Tokens, error) {
	const query = `
SELECT access_token, refresh_token
FROM channels_integrations_streamlabs
WHERE channel_id = $1 AND enabled = TRUE
LIMIT 1;
`

	var tokens provider.Tokens
	if err := c.tokenStoreDB.QueryRow(ctx, query, channelID).Scan(
		&tokens.AccessToken,
		&tokens.RefreshToken,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return provider.Tokens{}, streamlabsintegration.ErrNotFound
		}
		return provider.Tokens{}, err
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return provider.Tokens{}, errors.New("streamlabs integration has incomplete credentials")
	}
	return tokens, nil
}

func (c *Pgx) UpdateTokens(
	ctx context.Context,
	channelID string,
	tokens provider.Tokens,
) error {
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return errors.New("refuse to persist incomplete streamlabs credentials")
	}

	const query = `
UPDATE channels_integrations_streamlabs
SET
	"access_token" = $2,
	"refresh_token" = $3,
	updated_at = NOW()
WHERE channel_id = $1 AND enabled = TRUE
`

	command, err := c.tokenStoreDB.Exec(
		ctx,
		query,
		channelID,
		tokens.AccessToken,
		tokens.RefreshToken,
	)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return streamlabsintegration.ErrNotFound
	}
	return nil
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	model.StreamlabsIntegration,
	error,
) {
	query := `
SELECT id, channel_id, access_token, refresh_token, username, avatar, created_at, updated_at, enabled
FROM channels_integrations_streamlabs
WHERE channel_id = $1
LIMIT 1;
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, err
	}

	data, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.StreamlabsIntegration],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, streamlabsintegration.ErrNotFound
		}
		return model.Nil, err
	}

	return data, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	opts streamlabsintegration.UpdateOpts,
) error {
	query := `
UPDATE channels_integrations_streamlabs
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
		return streamlabsintegration.ErrNotFound
	}

	return err
}

func (c *Pgx) Delete(ctx context.Context, channelID string) error {
	query := `
DELETE FROM channels_integrations_streamlabs
WHERE channel_id = $1
`

	_, err := c.pool.Exec(ctx, query, channelID)
	return err
}

func (c *Pgx) Create(
	ctx context.Context,
	opts streamlabsintegration.CreateOpts,
) error {
	query := `
INSERT INTO channels_integrations_streamlabs (channel_id, enabled, access_token, refresh_token, username, avatar)
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
