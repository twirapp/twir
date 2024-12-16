package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/badges"
	"github.com/twirapp/twir/libs/repositories/badges/model"
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

var _ badges.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Badge, error) {
	query := `
SELECT id, name, enabled, created_at, file_name, ffz_slot
FROM badges
WHERE id = $1
LIMIT 1
`

	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, badges.ErrBadgeNotFound
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Badge])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetMany(ctx context.Context, input badges.GetManyInput) ([]model.Badge, error) {
	query := `
SELECT id, name, enabled, created_at, file_name, ffz_slot
FROM badges
WHERE enabled = $1
`

	rows, err := c.pool.Query(ctx, query, input.Enabled)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Badge])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE from badges
WHERE id = $1
`

	result, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return badges.ErrBadgeNotFound
	}

	return nil
}

func (c *Pgx) Create(ctx context.Context, input badges.CreateInput) (model.Badge, error) {
	query := `
INSERT INTO badges (id, name, enabled, ffz_slot, file_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, enabled, created_at, ffz_slot, file_name
`

	rows, err := c.pool.Query(
		ctx,
		query,
		input.ID,
		input.Name,
		input.Enabled,
		input.FFZSlot,
		input.FileName,
	)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Badge])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input badges.UpdateInput) (
	model.Badge,
	error,
) {
	updateBuilder := sq.Update("badges").
		Where("id = ?", id).
		Suffix("RETURNING id, name, enabled, created_at, ffz_slot, file_name")

	if input.Name != nil {
		updateBuilder = updateBuilder.Set("name", *input.Name)
	}

	if input.Enabled != nil {
		updateBuilder = updateBuilder.Set("enabled", *input.Enabled)
	}

	if input.FFZSlot != nil {
		updateBuilder = updateBuilder.Set("ffz_slot", *input.FFZSlot)
	}

	if input.FileName != nil {
		updateBuilder = updateBuilder.Set("file_name", *input.FileName)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Badge])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}
