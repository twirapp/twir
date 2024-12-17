package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/commands_response"
	"github.com/twirapp/twir/libs/repositories/commands_response/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ commands_response.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetManyByIDs(ctx context.Context, commandsIDs []uuid.UUID) (
	[]*model.Response,
	error,
) {
	if len(commandsIDs) == 0 {
		return nil, nil
	}

	query := `
SELECT id, "commandId", "order", text, twitch_category_id
FROM channels_commands_responses
WHERE "commandId" IN ($1)
`

	commandsIDsStrings := make([]string, len(commandsIDs))
	for i, id := range commandsIDs {
		commandsIDsStrings[i] = id.String()
	}

	rows, err := c.pool.Query(ctx, query, commandsIDsStrings)
	if err != nil {
		return nil, err
	}

	responses, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Response])
	if err != nil {
		return nil, err
	}

	result := make([]*model.Response, len(commandsIDs))
	for i, id := range commandsIDs {
		for _, r := range responses {
			if r.CommandID == id {
				result[i] = &r
				break
			}
		}
	}

	return result, nil
}
