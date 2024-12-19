package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories"
	"github.com/twirapp/twir/libs/repositories/greetings"
	"github.com/twirapp/twir/libs/repositories/greetings/model"
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

var _ greetings.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Greeting, error) {
	query := `
SELECT id, "channelId", "userId", enabled, text, "isReply", processed
FROM channels_greetings
WHERE "channelId" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Greeting])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Greeting, error) {
	query := `
SELECT id, "channelId", "userId", enabled, text, "isReply", processed
FROM channels_greetings
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.GreetingNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Greeting])
	if err != nil {
		return model.GreetingNil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input greetings.CreateInput) (model.Greeting, error) {
	query := `
INSERT INTO channels_greetings ("channelId", "userId", enabled, text, "isReply", processed)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, "channelId", "userId", enabled, text, "isReply", processed
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.UserID,
		input.Enabled,
		input.Text,
		input.IsReply,
		input.Processed,
	)
	if err != nil {
		return model.GreetingNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Greeting])
	if err != nil {
		return model.GreetingNil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input greetings.UpdateInput) (
	model.Greeting,
	error,
) {
	updateBuilder := sq.
		Update("channels_greetings").
		Where(squirrel.Eq{"id": id}).
		Suffix(`RETURNING id, "channelId", "userId", enabled, text, "isReply", processed`)
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
			`"userId"`:  input.UserID,
			"enabled":   input.Enabled,
			"text":      input.Text,
			`"isReply"`: input.IsReply,
			"processed": input.Processed,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.GreetingNil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.GreetingNil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Greeting])
	if err != nil {
		return model.GreetingNil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_greetings
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rows.RowsAffected())
	}

	return nil
}
