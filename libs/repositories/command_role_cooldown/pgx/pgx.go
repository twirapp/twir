package pgx

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/twirapp/twir/libs/entities/commandrolecooldownentity"
	"github.com/twirapp/twir/libs/repositories/command_role_cooldown"
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

var (
	_  command_role_cooldown.Repository = (*Pgx)(nil)
	sq                                  = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns = []string{
	"id",
	"command_id",
	"role_id",
	"cooldown",
	"created_at",
	"updated_at",
}

type scanModel struct {
	ID        uuid.UUID `db:"id"`
	CommandID uuid.UUID `db:"command_id"`
	RoleID    uuid.UUID `db:"role_id"`
	Cooldown  int       `db:"cooldown"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Pgx) GetByCommandID(ctx context.Context, commandID uuid.UUID) ([]entity.CommandRoleCooldown, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Select(selectColumns...).
		From("channels_commands_role_cooldowns").
		Where(squirrel.Eq{"command_id": commandID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByCommandID: failed to build select query: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetByCommandID: failed to execute select query: %w", err)
	}

	dbModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("GetByCommandID: failed to collect rows: %w", err)
	}

	result := make([]entity.CommandRoleCooldown, 0, len(dbModels))
	for _, dbModel := range dbModels {
		result = append(
			result, entity.CommandRoleCooldown{
				ID:        dbModel.ID,
				CommandID: dbModel.CommandID,
				RoleID:    dbModel.RoleID,
				Cooldown:  dbModel.Cooldown,
				CreatedAt: dbModel.CreatedAt,
				UpdatedAt: dbModel.UpdatedAt,
			},
		)
	}

	return result, nil
}

func (c *Pgx) GetByCommandIDs(ctx context.Context, commandIDs []uuid.UUID) ([]entity.CommandRoleCooldown, error) {
	if len(commandIDs) == 0 {
		return []entity.CommandRoleCooldown{}, nil
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Select(selectColumns...).
		From("channels_commands_role_cooldowns").
		Where(squirrel.Eq{"command_id": commandIDs}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetByCommandIDs: failed to build select query: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetByCommandIDs: failed to execute select query: %w", err)
	}

	dbModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("GetByCommandIDs: failed to collect rows: %w", err)
	}

	result := make([]entity.CommandRoleCooldown, 0, len(dbModels))
	for _, dbModel := range dbModels {
		result = append(
			result, entity.CommandRoleCooldown{
				ID:        dbModel.ID,
				CommandID: dbModel.CommandID,
				RoleID:    dbModel.RoleID,
				Cooldown:  dbModel.Cooldown,
				CreatedAt: dbModel.CreatedAt,
				UpdatedAt: dbModel.UpdatedAt,
			},
		)
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input command_role_cooldown.CreateInput) (entity.CommandRoleCooldown, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Insert("channels_commands_role_cooldowns").
		SetMap(
			map[string]any{
				"command_id": input.CommandID,
				"role_id":    input.RoleID,
				"cooldown":   input.Cooldown,
			},
		).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return entity.Nil, fmt.Errorf("Create: failed to build insert query: %w", err)
	}

	var id uuid.UUID
	var createdAt, updatedAt time.Time
	err = conn.QueryRow(ctx, query, args...).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return entity.Nil, fmt.Errorf("Create: failed to execute insert query: %w", err)
	}

	return entity.CommandRoleCooldown{
		ID:        id,
		CommandID: input.CommandID,
		RoleID:    input.RoleID,
		Cooldown:  input.Cooldown,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (c *Pgx) CreateMany(ctx context.Context, inputs []command_role_cooldown.CreateInput) error {
	if len(inputs) == 0 {
		return nil
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	insertBuilder := sq.Insert("channels_commands_role_cooldowns").
		Columns("command_id", "role_id", "cooldown")

	for _, input := range inputs {
		insertBuilder = insertBuilder.Values(input.CommandID, input.RoleID, input.Cooldown)
	}

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("CreateMany: failed to build insert query: %w", err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("CreateMany: failed to execute insert query: %w", err)
	}

	return nil
}

func (c *Pgx) DeleteByCommandID(ctx context.Context, commandID uuid.UUID) error {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Delete("channels_commands_role_cooldowns").
		Where(squirrel.Eq{"command_id": commandID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("DeleteByCommandID: failed to build delete query: %w", err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("DeleteByCommandID: failed to execute delete query: %w", err)
	}

	return nil
}

func (c *Pgx) DeleteByCommandIDAndRoleID(ctx context.Context, commandID, roleID uuid.UUID) error {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Delete("channels_commands_role_cooldowns").
		Where(
			squirrel.Eq{
				"command_id": commandID,
				"role_id":    roleID,
			},
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("DeleteByCommandIDAndRoleID: failed to build delete query: %w", err)
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("DeleteByCommandIDAndRoleID: failed to execute delete query: %w", err)
	}

	return nil
}
