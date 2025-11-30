package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelsintegrationsvalorant "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_valorant/model"
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

var _ channelsintegrationsvalorant.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) (model.ChannelIntegrationValorant, error) {
	query := `
SELECT id, channel_id, enabled, access_token, refresh_token, username, valorant_active_region, valorant_puuid
FROM channels_integrations_valorant
WHERE channel_id = $1
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationValorant],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	return result, nil
}

func (p Pgx) Create(
	ctx context.Context,
	input channelsintegrationsvalorant.CreateInput,
) (model.ChannelIntegrationValorant, error) {
	query := `
INSERT INTO channels_integrations_valorant (
	channel_id, enabled, access_token, refresh_token, username, valorant_active_region, valorant_puuid
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, channel_id, enabled, access_token, refresh_token, username, valorant_active_region, valorant_puuid
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)

	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Enabled,
		input.AccessToken,
		input.RefreshToken,
		input.UserName,
		input.ValorantActiveRegion,
		input.ValorantPuuid,
	)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationValorant],
	)
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (p Pgx) Update(
	ctx context.Context,
	id int,
	input channelsintegrationsvalorant.UpdateInput,
) error {
	builder := sq.Update("channels_integrations_valorant").Where(squirrel.Eq{"id": id})

	if input.Enabled != nil {
		builder = builder.Set("enabled", *input.Enabled)
	}

	if input.AccessToken != nil {
		builder = builder.Set("access_token", *input.AccessToken)
	}

	if input.RefreshToken != nil {
		builder = builder.Set("refresh_token", *input.RefreshToken)
	}

	if input.UserName != nil {
		builder = builder.Set("username", *input.UserName)
	}

	if input.ValorantActiveRegion != nil {
		builder = builder.Set("valorant_active_region", *input.ValorantActiveRegion)
	}

	if input.ValorantPuuid != nil {
		builder = builder.Set("valorant_puuid", *input.ValorantPuuid)
	}

	builder = builder.Set("updated_at", "now()")

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err = conn.Exec(ctx, query, args...)
	return err
}

func (p Pgx) Delete(
	ctx context.Context,
	id int,
) error {
	query := `
DELETE FROM channels_integrations_valorant
WHERE id = $1
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}
