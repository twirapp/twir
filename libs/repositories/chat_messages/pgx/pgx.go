package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
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

var _ chat_messages.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetMany(ctx context.Context, input chat_messages.GetManyInput) (
	[]model.ChatMessage,
	error,
) {
	perPage := input.PerPage
	if perPage == 0 {
		perPage = 20
	}

	if perPage > 1000 {
		perPage = 20
	}

	builder := sq.Select(
		"id",
		"channel_id",
		"user_id",
		"user_name",
		"user_display_name",
		"user_color",
		"text",
		"created_at",
		"updated_at",
	).From("chat_messages")

	if input.ChannelID != nil {
		builder = builder.Where(squirrel.Eq{"channel_id": *input.ChannelID})
	}

	if input.UserNameLike != nil && *input.UserNameLike != "" {
		builder = builder.Where(squirrel.ILike{"user_name": fmt.Sprintf("%%%s%%", *input.UserNameLike)})
	}

	if input.TextLike != nil && *input.TextLike != "" {
		builder = builder.Where(squirrel.ILike{"text": fmt.Sprintf("%%%s%%", *input.TextLike)})
	}

	builder = builder.OrderBy("created_at DESC").
		Offset(uint64(input.Page * perPage)).
		Limit(uint64(perPage))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChatMessage])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return models, nil
}

func (c *Pgx) Create(ctx context.Context, input chat_messages.CreateInput) error {
	query := `
INSERT INTO chat_messages (id, channel_id, user_id, text, user_name, user_display_name, user_color)
VALUES ($1, $2, $3, $4, $5, $6, &7);
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(
		ctx,
		query,
		input.ID,
		input.ChannelID,
		input.UserID,
		input.Text,
		input.UserName,
		input.UserDisplayName,
		input.UserColor,
	)
	return err
}

func (c *Pgx) CreateMany(ctx context.Context, inputs []chat_messages.CreateInput) error {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"chat_messages"},
		[]string{"id", "channel_id", "user_id", "text", "user_name", "user_display_name", "user_color"},
		pgx.CopyFromSlice(
			len(inputs),
			func(i int) ([]any, error) {
				return []any{
					inputs[i].ID,
					inputs[i].ChannelID,
					inputs[i].UserID,
					inputs[i].Text,
					inputs[i].UserName,
					inputs[i].UserDisplayName,
					inputs[i].UserColor,
				}, nil
			},
		),
	)
	if err != nil {
		return fmt.Errorf("failed to copy from: %w", err)
	}

	return nil
}
