package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/badges_users"
	"github.com/twirapp/twir/libs/repositories/badges_users/model"
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

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var _ badges_users.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetMany(ctx context.Context, input badges_users.GetManyInput) (
	[]model.BadgeUser,
	error,
) {
	query := `
SELECT id, badge_id, user_id, created_at
FROM badges_users
WHERE badge_id = $1
`

	rows, err := c.pool.Query(ctx, query, input.BadgeID)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.BadgeUser])
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input badges_users.CreateInput) (model.BadgeUser, error) {
	query := `
INSERT INTO "badges_users" (badge_id, user_id)
VALUES ($1, $2)
RETURNING id, badge_id, user_id, created_at
`

	rows, err := c.pool.Query(ctx, query, input.BadgeID, input.UserID)
	if err != nil {
		return model.BadgeUser{}, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.BadgeUser])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, input badges_users.DeleteInput) error {
	query := `
DELETE FROM badges_users
WHERE badge_id = $1 AND user_id = $2
`

	rows, err := c.pool.Exec(ctx, query, input.BadgeID, input.UserID)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return badges_users.ErrNotFound
	}

	return nil
}
