package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
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

var _ channelseventslist.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) CreateMany(ctx context.Context, inputs []channelseventslist.CreateInput) error {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"channels_events_list"},
		[]string{"channel_id", "user_id", "type", "data"},
		pgx.CopyFromSlice(
			len(inputs),
			func(i int) ([]any, error) {
				return []any{
					inputs[i].ChannelID,
					inputs[i].UserID,
					inputs[i].Type,
					inputs[i].Data,
				}, nil
			},
		),
	)
	if err != nil {
		return err
	}

	return nil
}
