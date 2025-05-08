package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/toxic_messages"
	"github.com/twirapp/twir/libs/repositories/toxic_messages/model"
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

func (c *Pgx) GetList(
	ctx context.Context,
	input toxic_messages.GetListInput,
) (toxic_messages.GetListOutput, error) {
	query := `
SELECT id, channel_id, reply_to_user_id, text, created_at
FROM toxic_messages
ORDER BY created_at DESC
LIMIT $1
OFFSET $2
`

	perPage := input.PerPage
	if perPage == 0 {
		perPage = 20
	}

	offset := input.Page * perPage

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, perPage, offset)
	if err != nil {
		return toxic_messages.GetListOutput{}, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ToxicMessage])
	if err != nil {
		return toxic_messages.GetListOutput{}, err
	}

	totalQuery := `
SELECT COUNT(*) FROM toxic_messages
`

	var total int
	err = conn.QueryRow(ctx, totalQuery).Scan(&total)
	if err != nil {
		return toxic_messages.GetListOutput{}, err
	}

	return toxic_messages.GetListOutput{
		Items: result,
		Total: total,
	}, nil
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
