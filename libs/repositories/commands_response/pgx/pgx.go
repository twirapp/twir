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
	"github.com/twirapp/twir/libs/repositories/commands_response"
	"github.com/twirapp/twir/libs/repositories/commands_response/model"
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

var _ commands_response.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) Create(ctx context.Context, input commands_response.CreateInput) (
	model.Response,
	error,
) {
	query := `
INSERT INTO channels_commands_responses("commandId", "order", text, twitch_category_id, online_only, offline_only)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, "commandId", "order", text, twitch_category_id, online_only, offline_only;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.CommandID,
		input.Order,
		input.Text,
		append([]string{}, input.TwitchCategoryIDs...),
		input.OnlineOnly,
		input.OfflineOnly,
	)
	if err != nil {
		return model.Nil, err
	}

	response, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Response])
	if err != nil {
		return model.Nil, err
	}

	return response, nil
}

func (c *Pgx) GetManyByIDs(ctx context.Context, commandsIDs []uuid.UUID) (
	[]model.Response,
	error,
) {
	if len(commandsIDs) == 0 {
		return nil, nil
	}

	query := `
SELECT id, "commandId", "order", text, twitch_category_id, online_only, offline_only
FROM channels_commands_responses
WHERE "commandId" = any($1);
`

	commandsIDsStrings := make([]string, len(commandsIDs))
	for i, id := range commandsIDs {
		commandsIDsStrings[i] = id.String()
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, commandsIDsStrings)
	if err != nil {
		return nil, err
	}

	responses, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Response])
	if err != nil {
		return nil, err
	}

	return responses, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_commands_responses
WHERE id = $1;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete: failed to execute delete query: %w", err)
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("delete: failed to delete exactly one row")
	}

	return nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input commands_response.UpdateInput,
) (model.Response, error) {
	updateBuilder := sq.Update("channels_commands_responses").
		Where(squirrel.Eq{"id": id}).
		Suffix(
			`RETURNING id, "commandId", "order", text, twitch_category_id`,
			"online_only",
			"offline_only",
		)
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]interface{}{
			"text":               input.Text,
			"order":              input.Order,
			"twitch_category_id": input.TwitchCategoryIDs,
			"online_only":        input.OnlineOnly,
			"offline_only":       input.OfflineOnly,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, fmt.Errorf("update: failed to build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("update: failed to execute query: %w", err)
	}

	response, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Response])
	if err != nil {
		return model.Nil, fmt.Errorf("update: failed to collect exactly one row: %w", err)
	}

	return response, nil
}
