package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/toxic_messages"
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

var _ toxic_messages.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) Create(ctx context.Context, input toxic_messages.CreateInput) error {
	if input.Text == "" || input.ChannelID == "" {
		return fmt.Errorf("text and channel_id are required")
	}

	query := `
INSERT INTO toxic_messages(channel_id, reply_to_user_id, text)
VALUES ($1, $2, $3)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, input.ChannelID, input.ReplyToUserID, input.Text)
	return err
}
