package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_categories_aliases"
	"github.com/twirapp/twir/libs/repositories/channels_categories_aliases/model"
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

var (
	_  channels_categories_aliases.Repository = (*Pgx)(nil)
	sq                                        = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetManyByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.ChannelCategoryAliase, error) {
	query := `
SELECT id, channel_id, alias, category_id
FROM channels_categories_aliases
WHERE channel_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	aliases, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelCategoryAliase])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return aliases, nil
}

func (c *Pgx) Create(ctx context.Context, input channels_categories_aliases.CreateInput) error {
	query := `
INSERT INTO channels_categories_aliases (channel_id, alias, category_id)
VALUES ($1, $2, $3)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, input.ChannelID, input.Alias, input.CategoryID)
	return err
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_categories_aliases
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id.String())
	return err
}
