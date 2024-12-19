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
	"github.com/twirapp/twir/libs/repositories/roles"
	"github.com/twirapp/twir/libs/repositories/roles/model"
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

var _ roles.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Role, error) {
	query := `
SELECT id, "channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time
FROM channels_roles
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.RoleNil, fmt.Errorf("GetByID: failed to execute select query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return model.RoleNil, fmt.Errorf("GetByID: failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) GetManyByIDS(ctx context.Context, ids []uuid.UUID) ([]model.Role, error) {
	query := `
SELECT id, "channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time
FROM channels_roles
WHERE id = ANY($1)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("GetManyByIDS: failed to execute select query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return nil, fmt.Errorf("GetManyByIDS: failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Role, error) {
	query := `
SELECT id, "channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time
FROM channels_roles
WHERE "channelId" = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("GetManyByChannelID: failed to execute select query: %w", err)
	}

	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return nil, fmt.Errorf("GetManyByChannelID: failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input roles.CreateInput) (model.Role, error) {
	query := `
INSERT INTO channels_roles("channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, "channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Name,
		input.Type,
		input.Permissions,
		input.RequiredMessages,
		input.RequiredUsedChannelPoints,
		input.RequiredWatchTime,
	)
	if err != nil {
		return model.RoleNil, fmt.Errorf("cannot create role: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return model.RoleNil, fmt.Errorf("cannot create role: failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input roles.UpdateInput) (
	model.Role,
	error,
) {
	updateBuilder := sq.
		Update("channels_roles").
		Where(squirrel.Eq{"id": id}).
		Suffix(`RETURNING id, "channelId", name, type, permissions, required_messages, required_used_channel_points, required_watch_time`)
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]any{
			"name":                         input.Name,
			"permissions":                  input.Permissions,
			"required_messages":            input.RequiredMessages,
			"required_used_channel_points": input.RequiredUsedChannelPoints,
			"required_watch_time":          input.RequiredWatchTime,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.RoleNil, fmt.Errorf("cannot update role: failed to build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.RoleNil, fmt.Errorf("cannot update role: failed to execute query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Role])
	if err != nil {
		return model.RoleNil, fmt.Errorf("cannot update role: failed to collect rows: %w", err)
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_roles
WHERE id = $1 AND type = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, id, model.ChannelRoleTypeCustom)
	if err != nil {
		return fmt.Errorf("cannot delete role: %w", err)
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("cannot delete role: role not found")
	}

	return nil
}
