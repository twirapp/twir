package pgx

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/roles_users"
	"github.com/twirapp/twir/libs/repositories/roles_users/model"
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

var _ roles_users.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetManyByRoleID(ctx context.Context, roleID uuid.UUID) ([]model.RoleUser, error) {
	query := `
SELECT id, "userId", "roleId"
FROM channels_roles_users
WHERE "roleId" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.RoleUser])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input roles_users.CreateInput) (model.RoleUser, error) {
	query := `
INSERT INTO channels_roles_users("userId", "roleId")
VALUES ($1, $2)
RETURNING id, "userId", "roleId"
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.UserID, input.RoleID)
	if err != nil {
		return model.RoleUserNil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.RoleUser])
	if err != nil {
		return model.RoleUserNil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_roles_users
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("role not found")
	}

	return nil
}

func (c *Pgx) DeleteManyByRoleID(ctx context.Context, roleID uuid.UUID) error {
	query := `
DELETE FROM channels_roles_users
WHERE "roleId" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, roleID)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	return nil
}

func (c *Pgx) CreateMany(ctx context.Context, inputs []roles_users.CreateInput) (
	[]model.RoleUser,
	error,
) {
	insertBuilder := sq.Insert("channels_roles_users").
		Columns(`"userId"`, `"roleId"`).
		Suffix(`RETURNING id, "userId", "roleId"`)
	for _, input := range inputs {
		insertBuilder = insertBuilder.Values(input.UserID, input.RoleID)
	}

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build insert query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute insert query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.RoleUser])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return result, nil
}
