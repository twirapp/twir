package pgx

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelsintegrationsspotify "github.com/twirapp/twir/libs/repositories/channels_integrations_spotify"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_spotify/model"
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

var _ channelsintegrationsspotify.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	model.ChannelIntegrationSpotify,
	error,
) {
	query := `
SELECT id, access_token, refresh_token, enabled, scopes, channel_id, avatar_uri, username, created_at, updated_at
FROM channels_integrations_spotify
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationSpotify],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}

		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input channelsintegrationsspotify.UpdateInput,
) error {
	updateBuilder := sq.Update("channels_integrations_spotify").Where(squirrel.Eq{"id": id})

	if input.AccessToken != nil {
		updateBuilder = updateBuilder.Set("access_token", *input.AccessToken)
	}

	if input.RefreshToken != nil {
		updateBuilder = updateBuilder.Set("refresh_token", *input.RefreshToken)
	}

	if input.AvatarURI != nil {
		updateBuilder = updateBuilder.Set("avatar_uri", *input.AvatarURI)
	}

	if input.Username != nil {
		updateBuilder = updateBuilder.Set("username", *input.Username)
	}

	if input.Scopes != nil {
		updateBuilder = updateBuilder.Set("scopes", *input.Scopes)
	}

	updateBuilder = updateBuilder.Set("updated_at", squirrel.Expr("NOW()"))

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Create(
	ctx context.Context,
	input channelsintegrationsspotify.CreateInput,
) (model.ChannelIntegrationSpotify, error) {
	query := `
INSERT INTO channels_integrations_spotify (access_token, refresh_token, avatar_uri, username, scopes, channel_id, enabled)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, access_token, refresh_token, enabled, scopes, channel_id, avatar_uri, username, created_at, updated_at
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.AccessToken,
		input.RefreshToken,
		input.AvatarURI,
		input.Username,
		input.Scopes,
		input.ChannelID,
		true,
	)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationSpotify],
	)
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_integrations_spotify
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}
