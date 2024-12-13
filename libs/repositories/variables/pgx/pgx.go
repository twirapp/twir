package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/variables"
	"github.com/twirapp/twir/libs/repositories/variables/model"
)

type Opts struct {
	Pgx *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.Pgx,
	}
}

func NewFx(pgxpool *pgxpool.Pool) *Pgx {
	return New(
		Opts{
			Pgx: pgxpool,
		},
	)
}

var _ variables.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) (
	[]model.CustomVariable,
	error,
) {
	query := `
SELECT id, "channelId", description, "evalValue", name, response, type
FROM channels_customvars
WHERE "channelId" = $1
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	vars, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.CustomVariable])
	if err != nil {
		return nil, err
	}

	return vars, nil
}

func (c *Pgx) CountByChannelID(ctx context.Context, channelID string) (int, error) {
	query := `
SELECT COUNT(*)
FROM channels_customvars
WHERE "channelId" = $1
`

	var count int
	err := c.pool.QueryRow(ctx, query, channelID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.CustomVariable, error) {
	query := `
SELECT id, "channelId", description, "evalValue", name, response, type
FROM channels_customvars
WHERE id = $1
LIMIT 1
`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	variable, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CustomVariable])
	if err != nil {
		return model.Nil, err
	}

	return variable, nil
}

func (c *Pgx) Create(ctx context.Context, input variables.CreateInput) (
	model.CustomVariable,
	error,
) {
	query := `
INSERT INTO channels_customvars ("channelId", description, "evalValue", name, response, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, "channelId", description, "evalValue", name, response, type
`

	rows, err := c.pool.Query(
		ctx,
		query,
		input.ChannelID,
		input.Description,
		input.EvalValue,
		input.Name,
		input.Response,
		input.Type,
	)
	if err != nil {
		return model.Nil, err
	}

	variable, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.CustomVariable])
	if err != nil {
		return model.Nil, err
	}

	return variable, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input variables.UpdateInput,
) (model.CustomVariable, error) {
	updateBuilder := sq.Update("channels_customvars")

	if input.Description != nil {
		updateBuilder = updateBuilder.Set("description", *input.Description)
	}

	if input.EvalValue != nil {
		updateBuilder = updateBuilder.Set(`"evalValue"`, *input.EvalValue)
	}

	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", *input.Name)
	}

	if input.Response != nil {
		updateBuilder = updateBuilder.Set("response", *input.Response)
	}

	if input.Type != nil {
		updateBuilder = updateBuilder.Set("type", *input.Type)
	}

	updateBuilder = updateBuilder.Where(squirrel.Eq{"id": id})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	_, err = c.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_customvars
WHERE id = $1
`

	rows, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return variables.ErrNotFound
	}

	return nil
}
