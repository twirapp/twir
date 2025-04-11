package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_info_history"
	"github.com/twirapp/twir/libs/repositories/channels_info_history/model"
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

var _ channels_info_history.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p *Pgx) GetMany(
	ctx context.Context,
	input channels_info_history.GetManyInput,
) ([]model.ChannelInfoHistory, error) {
	limit := input.Limit
	if limit > 100 || limit == 0 {
		limit = 10
	}

	if input.ChannelID == "" {
		return nil, fmt.Errorf("channel id cannot be empty")
	}

	builder := sq.Select(
		`id`, `"channelId"`, `"createdAt"`, `title`, `category`,
	).
		From("channels_info_history").
		Where(`"channelId" = ?`, input.ChannelID).
		Limit(uint64(limit))

	if input.UniqueBy != nil {
		switch *input.UniqueBy {
		case channels_info_history.UniqueByCategory:
			builder = builder.
				Options("DISTINCT ON (category)").
				OrderBy(`category, "createdAt" DESC`)
		case channels_info_history.UniqueByTitle:
			builder = builder.
				Options("DISTINCT ON (title)").
				OrderBy(`title, "createdAt" DESC`)
		}
	} else {
		builder = builder.OrderBy(`"createdAt" DESC`)
	}

	if !input.After.IsZero() {
		builder = builder.Where(`"createdAt" >= ?`, input.After)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build error: %w", err)
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelInfoHistory])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("collect error: %w", err)
	}

	return result, nil
}

func (p *Pgx) Create(ctx context.Context, input channels_info_history.CreateInput) error {
	query := `
INSERT INTO channels_info_history ("channelId", title, category)
VALUES ($1, $2, $3)
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, input.ChannelID, input.Title, input.Category)
	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}

	return nil
}
