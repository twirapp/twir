package pgx

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	commandsrepositorypgx "github.com/twirapp/twir/libs/repositories/commands/pgx"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
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

var _ commands_with_groups_and_responses.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns []string

func init() {
	var columns []string
	for _, column := range commandsrepositorypgx.SelectColumns {
		columns = append(columns, "c."+column)
	}

	columns = append(
		columns,
		`
COALESCE(json_agg(json_build_object(
	'id', r.id,
	'text', r.text,
	'commandId', r."commandId",
	'order', r."order",
	'twitch_category_id', r."twitch_category_id"
)) FILTER (WHERE r.id IS NOT NULL), '[]'::json) as responses
`,
		`
json_build_object(
	'id', g.id,
	'channelId', g."channelId",
	'name', g."name",
	'color', g.color
) as group
`,
	)

	selectColumns = append(selectColumns, columns...)
}

func (c *Pgx) GetManyByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.CommandWithGroupAndResponses, error) {
	selectBuilder := sq.Select(selectColumns...).
		From("channels_commands c").
		LeftJoin(`channels_commands_groups g ON c."groupId" = g.id`).
		LeftJoin(`channels_commands_responses r ON c.id = r."commandId"`).
		Where(`c."channelId" = ?`, channelID).
		GroupBy("c.id", "g.id")

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("responses GetManyByChannelID: failed to execute select query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CommandWithGroupAndResponses])
	if err != nil {
		return nil, fmt.Errorf("responses GetManyByChannelID: failed to collect rows: %w", err)
	}

	slices.SortFunc(
		result, func(i, j model.CommandWithGroupAndResponses) int {
			return cmp.Compare(i.Command.ID.String(), j.Command.ID.String())
		},
	)

	return result, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (
	model.CommandWithGroupAndResponses,
	error,
) {
	selectBuilder := sq.Select(selectColumns...).
		From("channels_commands c").
		LeftJoin(`channels_commands_groups g ON c."groupId" = g.id`).
		LeftJoin(`channels_commands_responses r ON c.id = r."commandId"`).
		Where(`c.id = ?`, id).
		GroupBy("c.id", "g.id")

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("responses GetByID: failed to execute select query: %w", err)
	}

	command, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.CommandWithGroupAndResponses],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, commands_with_groups_and_responses.ErrNotFound
		}
		return model.Nil, fmt.Errorf("responses GetByID: failed to collect rows: %w", err)
	}

	return command, nil
}
