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

func (c *Pgx) CountBy(ctx context.Context, input channelseventslist.CountByInput) (int64, error) {
	query := sq.Select("COUNT(*)").From("channels_events_list")

	if input.ChannelID != nil {
		query = query.Where(squirrel.Eq{"channel_id": *input.ChannelID})
	}
	if input.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *input.UserID})
	}
	if input.Type != nil {
		query = query.Where(squirrel.Eq{"type": *input.Type})
	}
	if input.CreatedAtGTE != nil {
		query = query.Where(squirrel.GtOrEq{"created_at": *input.CreatedAtGTE})
	}
	if input.CreatedAtLTE != nil {
		query = query.Where(squirrel.LtOrEq{"created_at": *input.CreatedAtLTE})
	}
	if input.CreatedAtGT != nil {
		query = query.Where(squirrel.Gt{"created_at": *input.CreatedAtGT})
	}
	if input.CreatedAtLT != nil {
		query = query.Where(squirrel.Lt{"created_at": *input.CreatedAtLT})
	}
	if input.CreatedAtEQ != nil {
		query = query.Where(squirrel.Eq{"created_at": *input.CreatedAtEQ})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	err = conn.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) Create(ctx context.Context, input channelseventslist.CreateInput) error {
	query := `
INSERT INTO channels_events_list (channel_id, user_id, type, data)
VALUES ($1, $2, $3, $4)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, input.ChannelID, input.UserID, input.Type, input.Data)
	return err
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
