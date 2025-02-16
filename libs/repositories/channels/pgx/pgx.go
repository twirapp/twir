package pgx

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels/model"
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

var _ channels.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByID(ctx context.Context, channelID string) (model.Channel, error) {
	query := `
SELECT "id", "isEnabled", "isTwitchBanned", "isBotMod", "botId"
FROM channels
WHERE "id" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels.ErrNotFound
		}
	}

	return result, nil
}

func (c *Pgx) GetMany(ctx context.Context, input channels.GetManyInput) ([]model.Channel, error) {
	selectBuilder := sq.
		Select(
			"id",
			`"isEnabled"`,
			`"isTwitchBanned"`,
			`"isBotMod"`,
			`"botId"`,
		).
		From("channels")

	if input.Enabled != nil {
		selectBuilder = selectBuilder.Where(`"isEnabled" = ?`, *input.Enabled)
	}

	// not need to use defaults because i wanna select all channels
	if input.PerPage > 0 {
		selectBuilder = selectBuilder.Limit(uint64(input.PerPage))
	}

	// not need to use defaults because i wanna select all channels
	if input.Page > 0 {
		selectBuilder = selectBuilder.Offset(uint64(input.Page * input.PerPage))
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Channel])
	if err != nil {
		return nil, err
	}

	return result, nil
}
