package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelsintegrationslastfm "github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm/model"
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

var _ channelsintegrationslastfm.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) (model.ChannelIntegrationLastfm, error) {
	query := `
SELECT id, channel_id, enabled, session_key, username, avatar
FROM channels_integrations_lastfm
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
		pgx.RowToStructByName[model.ChannelIntegrationLastfm],
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
	input channelsintegrationslastfm.CreateInput,
) (model.ChannelIntegrationLastfm, error) {
	query := `
INSERT INTO channels_integrations_lastfm (
	channel_id, enabled, session_key, username, avatar
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, channel_id, enabled, session_key, username, avatar
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)

	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Enabled,
		input.SessionKey,
		input.UserName,
		input.Avatar,
	)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.ChannelIntegrationLastfm],
	)
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (p Pgx) Update(
	ctx context.Context,
	id int,
	input channelsintegrationslastfm.UpdateInput,
) error {
	builder := sq.Update("channels_integrations_lastfm").Where(squirrel.Eq{"id": id})

	if input.Enabled != nil {
		builder = builder.Set("enabled", *input.Enabled)
	}

	if input.SessionKey != nil {
		builder = builder.Set("session_key", *input.SessionKey)
	}

	if input.UserName != nil {
		builder = builder.Set("username", *input.UserName)
	}

	if input.Avatar != nil {
		builder = builder.Set("avatar", *input.Avatar)
	}

	builder = builder.Set("updated_at", squirrel.Expr("NOW()"))

	query, args, err := builder.ToSql()
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
DELETE FROM channels_integrations_lastfm
WHERE id = $1
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}
