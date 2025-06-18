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
	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands/model"
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

var _ commands.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Command, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Select(SelectColumns...).
		From("channels_commands").
		Where(squirrel.Eq{`"channelId"`: channelID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetManyByChannelID: failed to build select query: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("GetManyByChannelID: failed to execute select query: %w", err)
	}

	cmd, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Command])
	if err != nil {
		return nil, fmt.Errorf("GetManyByChannelID: failed to collect rows: %w", err)
	}

	return cmd, nil
}

var SelectColumns = []string{
	"id",
	"name",
	"cooldown",
	`"cooldownType"`,
	"enabled",
	"aliases",
	"description",
	"visible",
	`"channelId"`,
	`"default"`,
	`"defaultName"`,
	`"module"`,
	`is_reply`,
	`"keepResponsesOrder"`,
	`"deniedUsersIds"`,
	`"allowedUsersIds"`,
	`"rolesIds"`,
	`online_only`,
	`offline_only`,
	`cooldown_roles_ids`,
	`enabled_categories`,
	`"requiredWatchTime"`,
	`"requiredMessages"`,
	`"requiredUsedChannelPoints"`,
	`"groupId"`,
	`expires_at`,
	`expires_type`,
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Command, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Select(SelectColumns...).
		From("channels_commands").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return model.Nil, fmt.Errorf("GetByID: failed to build select query: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("GetByID: failed to execute select query: %w", err)
	}

	command, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Command])
	if err != nil {
		return model.Nil, fmt.Errorf("GetByID: failed to collect exactly one row: %w", err)
	}

	return command, nil
}

func (c *Pgx) Create(ctx context.Context, input commands.CreateInput) (model.Command, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query, args, err := sq.Insert("channels_commands").
		SetMap(
			map[string]any{
				`"channelId"`:                 input.ChannelID,
				"name":                        input.Name,
				"cooldown":                    input.Cooldown,
				`"default"`:                   false,
				`"defaultName"`:               nil,
				"module":                      "CUSTOM",
				`"cooldownType"`:              input.CooldownType,
				"enabled":                     input.Enabled,
				"aliases":                     append([]string{}, input.Aliases...),
				"description":                 input.Description,
				"visible":                     input.Visible,
				"is_reply":                    input.IsReply,
				`"keepResponsesOrder"`:        input.KeepResponsesOrder,
				`"deniedUsersIds"`:            append([]string{}, input.DeniedUsersIDS...),
				`"allowedUsersIds"`:           append([]string{}, input.AllowedUsersIDS...),
				`"rolesIds"`:                  append([]string{}, input.RolesIDS...),
				"online_only":                 input.OnlineOnly,
				"offline_only":                input.OfflineOnly,
				`"cooldown_roles_ids"`:        append([]string{}, input.CooldownRolesIDs...),
				`"enabled_categories"`:        input.EnabledCategories,
				`"requiredWatchTime"`:         input.RequiredWatchTime,
				`"requiredMessages"`:          input.RequiredMessages,
				`"requiredUsedChannelPoints"`: input.RequiredUsedChannelPoints,
				`"groupId"`:                   input.GroupID,
				`"expires_at"`:                input.ExpiresAt,
				`"expires_type"`:              input.ExpiresType,
			},
		).Suffix("RETURNING ID").ToSql()
	if err != nil {
		return model.Nil, fmt.Errorf("insert: failed to build query: %w", err)
	}

	rows := conn.QueryRow(ctx, query, args...)
	var id uuid.UUID
	if err := rows.Scan(&id); err != nil {
		return model.Nil, fmt.Errorf("insert: failed to scan id: %w", err)
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM "channels_commands"
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	rows, err := conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete: failed to execute delete query: %w", err)
	}

	if rows.RowsAffected() != 1 {
		return fmt.Errorf("delete: command not found")
	}

	return nil
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input commands.UpdateInput) (
	model.Command,
	error,
) {
	updateBuilder := sq.Update("channels_commands").
		Where(squirrel.Eq{"id": id})
	updateBuilder = repositories.SquirrelApplyPatch(
		updateBuilder,
		map[string]interface{}{
			"name":                        input.Name,
			"cooldown":                    input.Cooldown,
			`"cooldownType"`:              input.CooldownType,
			"enabled":                     input.Enabled,
			"aliases":                     input.Aliases,
			"description":                 input.Description,
			"visible":                     input.Visible,
			"is_reply":                    input.IsReply,
			`"keepResponsesOrder"`:        input.KeepResponsesOrder,
			`"deniedUsersIds"`:            input.DeniedUsersIDS,
			`"allowedUsersIds"`:           input.AllowedUsersIDS,
			`"rolesIds"`:                  input.RolesIDS,
			"online_only":                 input.OnlineOnly,
			`"offline_only"`:              input.OfflineOnly,
			`"cooldown_roles_ids"`:        input.CooldownRolesIDs,
			`"enabled_categories"`:        input.EnabledCategories,
			`"requiredWatchTime"`:         input.RequiredWatchTime,
			`"requiredMessages"`:          input.RequiredMessages,
			`"requiredUsedChannelPoints"`: input.RequiredUsedChannelPoints,
			`"groupId"`:                   input.GroupID,
			`"expires_at"`:                input.ExpiresAt,
			`"expires_type"`:              input.ExpiresType,
		},
	)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, fmt.Errorf("update: failed to build query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("update: failed to execute query: %w", err)
	}

	if rows.RowsAffected() != 1 {
		return model.Nil, fmt.Errorf("update: command not found")
	}

	return c.GetByID(ctx, id)
}
