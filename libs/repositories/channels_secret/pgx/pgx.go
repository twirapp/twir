package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_secret"
	"github.com/twirapp/twir/libs/repositories/channels_secret/model"
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

var (
	_  channels_secret.Repository = (*Pgx)(nil)
	sq                            = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) (
	[]model.ChannelSecret,
	error,
) {
	query := `
SELECT id, channel_id, name, description, value, created_at, updated_at
FROM channels_secrets
WHERE channel_id = $1
ORDER BY name
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	secrets, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ChannelSecret])
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.ChannelSecret, error) {
	query := `
SELECT id, channel_id, name, description, value, created_at, updated_at
FROM channels_secrets
WHERE id = $1
LIMIT 1
`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	secret, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelSecret])
	if err != nil {
		return model.Nil, err
	}

	return secret, nil
}

func (c *Pgx) Create(ctx context.Context, input channels_secret.CreateInput) (
	model.ChannelSecret,
	error,
) {
	query := `
INSERT INTO channels_secrets (channel_id, name, description, value)
VALUES ($1, $2, $3, $4)
RETURNING id, channel_id, name, description, value, created_at, updated_at
`

	rows, err := c.pool.Query(
		ctx,
		query,
		input.ChannelID,
		input.Name,
		input.Description,
		input.Value,
	)
	if err != nil {
		return model.Nil, err
	}

	secret, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.ChannelSecret])
	if err != nil {
		return model.Nil, err
	}

	return secret, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input channels_secret.UpdateInput,
) (model.ChannelSecret, error) {
	updateBuilder := sq.Update("channels_secrets")

	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", *input.Name)
	}

	if input.Description != nil {
		updateBuilder = updateBuilder.Set("description", *input.Description)
	}

	if input.Value != nil {
		updateBuilder = updateBuilder.Set("value", *input.Value)
	}

	updateBuilder = updateBuilder.Set("updated_at", squirrel.Expr("NOW()"))
	updateBuilder = updateBuilder.Where(squirrel.Eq{"id": id})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	_, err = c.pool.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_secrets
WHERE id = $1
`

	tag, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() != 1 {
		return channels_secret.ErrNotFound
	}

	return nil
}
