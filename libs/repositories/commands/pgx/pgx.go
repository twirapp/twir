package pgx

import (
	"cmp"
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/commands"
	"github.com/twirapp/twir/libs/repositories/commands/model"
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

var _ commands.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

var selectColumns = []string{
	"c.id",
	"c.name",
	"c.cooldown",
	`c."cooldownType"`,
	"c.enabled",
	"c.aliases",
	"c.description",
	"c.visible",
	`c."channelId"`,
	`c."default"`,
	`c."defaultName"`,
	`c."module"`,
	`c.is_reply`,
	`c."keepResponsesOrder"`,
	`c."deniedUsersIds"`,
	`c."allowedUsersIds"`,
	`c."rolesIds"`,
	`c.online_only`,
	`c.cooldown_roles_ids`,
	`c.enabled_categories`,
	`c."requiredWatchTime"`,
	`c."requiredMessages"`,
	`c."requiredUsedChannelPoints"`,
	`c."groupId"`,
	`c.expires_at`,
	`c.expires_type`,
	`g.id group_id`,
	`g."channelId" group_channel_id`,
	`g."name" group_name`,
	`g.color group_color`,
	"r.id response_id",
	"r.text response_text",
	`r."commandId" response_command_id`,
	"r.order response_order",
	"r.twitch_category_id response_twitch_category_id",
}

func (c *Pgx) scanRow(rows pgx.Rows) (model.Command, error) {
	var command model.Command
	var commandDefaultName, commandDescription, commandGroupID sql.Null[string]
	var commandCooldown sql.Null[int]
	var commandExpiresAt sql.Null[time.Time]
	var commandExpiresType sql.Null[model.ExpireType]

	var responseID, responseCommandID sql.Null[uuid.UUID]
	var responseText sql.Null[string]
	var responseTwitchCategoryID []string
	var responseOrder sql.Null[int]

	var groupId sql.Null[uuid.UUID]
	var groupChannelID, groupName, groupColor sql.Null[string]

	if err := rows.Scan(
		&command.ID,
		&command.Name,
		&commandCooldown,
		&command.CooldownType,
		&command.Enabled,
		&command.Aliases,
		&commandDescription,
		&command.Visible,
		&command.ChannelID,
		&command.Default,
		&commandDefaultName,
		&command.Module,
		&command.IsReply,
		&command.KeepResponsesOrder,
		&command.DeniedUsersIDS,
		&command.AllowedUsersIDS,
		&command.RolesIDS,
		&command.OnlineOnly,
		&command.CooldownRolesIDs,
		&command.EnabledCategories,
		&command.RequiredWatchTime,
		&command.RequiredMessages,
		&command.RequiredUsedChannelPoints,
		&commandGroupID,
		&commandExpiresAt,
		&commandExpiresType,
		&groupId,
		&groupChannelID,
		&groupName,
		&groupColor,
		&responseID,
		&responseText,
		&responseCommandID,
		&responseOrder,
		&responseTwitchCategoryID,
	); err != nil {
		return model.Nil, err
	}

	if commandDefaultName.Valid {
		command.DefaultName = &commandDefaultName.V
	}

	if commandDescription.Valid {
		command.Description = &commandDescription.V
	}

	if commandCooldown.Valid {
		command.Cooldown = &commandCooldown.V
	}

	if commandExpiresAt.Valid {
		command.ExpiresAt = &commandExpiresAt.V
	}

	if commandExpiresType.Valid {
		command.ExpiresType = &commandExpiresType.V
	}

	if responseID.Valid {
		command.Responses = append(
			command.Responses,
			model.Response{
				ID:                responseID.V,
				Text:              &responseText.V,
				CommandID:         responseCommandID.V,
				Order:             responseOrder.V,
				TwitchCategoryIDs: responseTwitchCategoryID,
			},
		)
	}

	if groupId.Valid {
		command.GroupID = &groupId.V
		command.Group = &model.Group{
			ID:        groupId.V,
			ChannelID: groupChannelID.V,
			Name:      groupName.V,
			Color:     groupColor.V,
		}
	}

	return command, nil
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Command, error) {
	selectBuilder := sq.Select(selectColumns...).
		From("channels_commands c").
		LeftJoin(`channels_commands_groups g ON c."groupId" = g.id`).
		LeftJoin(`channels_commands_responses r ON c.id = r."commandId"`).
		Where(`c."channelId" = ?`, channelID)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	commandsMap := make(map[uuid.UUID]*model.Command)
	for rows.Next() {
		command, err := c.scanRow(rows)
		if err != nil {
			return nil, err
		}

		if _, ok := commandsMap[command.ID]; !ok {
			commandsMap[command.ID] = &command
		}

		for _, r := range command.Responses {
			for _, responses := range commandsMap[command.ID].Responses {
				if responses.ID == r.ID {
					continue
				}

				commandsMap[command.ID].Responses = append(
					commandsMap[command.ID].Responses,
					r,
				)
			}
		}
	}

	result := make([]model.Command, 0, len(commandsMap))
	for _, cmd := range commandsMap {
		result = append(result, *cmd)
	}

	slices.SortFunc(
		result, func(i, j model.Command) int {
			return cmp.Compare(i.ID.String(), j.ID.String())
		},
	)

	return result, nil
}
