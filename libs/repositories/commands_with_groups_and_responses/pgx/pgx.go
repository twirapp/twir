package pgx

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	commandmodel "github.com/twirapp/twir/libs/repositories/commands/model"
	commandsrepositorypgx "github.com/twirapp/twir/libs/repositories/commands/pgx"
	groupmodel "github.com/twirapp/twir/libs/repositories/commands_group/model"
	responsemodel "github.com/twirapp/twir/libs/repositories/commands_response/model"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
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

var _ commands_with_groups_and_responses.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns []string

func init() {
	var columns []string
	for _, column := range commandsrepositorypgx.SelectColumns {
		columns = append(columns, "c."+column)
	}

	columns = append(
		columns,
		`g.id group_id`,
		`g."channelId" group_channel_id`,
		`g."name" group_name`,
		`g.color group_color`,
		"r.id response_id",
		"r.text response_text",
		`r."commandId" response_command_id`,
		"r.order response_order",
		"r.twitch_category_id response_twitch_category_id",
		"r.online_only response_online_only",
		"r.offline_only response_offline_only",
	)

	selectColumns = append(selectColumns, columns...)
}

type scanModel struct {
	Command  commandmodel.Command
	Group    *groupmodel.Group
	Response *responsemodel.Response
}

func (c *Pgx) scanRow(rows pgx.Rows) (scanModel, error) {
	var command commandmodel.Command
	var group groupmodel.Group

	var commandDefaultName, commandDescription, commandGroupID sql.Null[string]
	var commandCooldown sql.Null[int]
	var commandExpiresAt sql.Null[time.Time]
	var commandExpiresType sql.Null[commandmodel.ExpireType]

	var responseID, responseCommandID sql.Null[uuid.UUID]
	var responseText sql.Null[string]
	var responseTwitchCategoryID []string
	var responseOrder sql.Null[int]
	var responseOnlineOnly, responseOfflineOnly sql.Null[bool]

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
		&command.OfflineOnly,
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
		&responseOnlineOnly,
		&responseOfflineOnly,
	); err != nil {
		return scanModel{}, fmt.Errorf("responses failed to scan row: %w", err)
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

	var response *responsemodel.Response
	if responseID.Valid {
		response = &responsemodel.Response{
			ID:                responseID.V,
			Text:              &responseText.V,
			CommandID:         responseCommandID.V,
			Order:             responseOrder.V,
			TwitchCategoryIDs: responseTwitchCategoryID,
		}
	}

	if groupId.Valid {
		command.GroupID = &groupId.V
		group = groupmodel.Group{
			ID:        groupId.V,
			ChannelID: groupChannelID.V,
			Name:      groupName.V,
			Color:     groupColor.V,
		}
	}

	return scanModel{
		Command:  command,
		Group:    &group,
		Response: response,
	}, nil
}

func (c *Pgx) GetManyByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.CommandWithGroupAndResponses, error) {
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
		return nil, fmt.Errorf("responses GetManyByChannelID: failed to execute select query: %w", err)
	}

	commandsMap := make(map[uuid.UUID]*model.CommandWithGroupAndResponses)
	for rows.Next() {
		command, err := c.scanRow(rows)
		if err != nil {
			return nil, err
		}

		if _, ok := commandsMap[command.Command.ID]; !ok {
			commandsMap[command.Command.ID] = &model.CommandWithGroupAndResponses{
				Command:   command.Command,
				Group:     command.Group,
				Responses: []responsemodel.Response{},
			}
		}

		if command.Response != nil {
			commandsMap[command.Command.ID].Responses = append(
				commandsMap[command.Command.ID].Responses,
				*command.Response,
			)
		}
	}

	result := make([]model.CommandWithGroupAndResponses, 0, len(commandsMap))
	for _, cmd := range commandsMap {
		result = append(result, *cmd)
	}

	slices.SortFunc(
		result, func(i, j model.CommandWithGroupAndResponses) int {
			return cmp.Compare(i.Command.ID.String(), j.Command.ID.String())
		},
	)

	return result, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (
	model.CommandWithGroupAndResponses,
	error,
) {
	selectBuilder := sq.Select(selectColumns...).
		From("channels_commands c").
		LeftJoin(`channels_commands_groups g ON c."groupId" = g.id`).
		LeftJoin(`channels_commands_responses r ON c.id = r."commandId"`).
		Where(`c.id = ?`, id)

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("responses GetByID: failed to execute select query: %w", err)
	}

	defer rows.Close()

	var command model.CommandWithGroupAndResponses
	for rows.Next() {
		cmd, err := c.scanRow(rows)
		if err != nil {
			return model.Nil, err
		}

		if command.Command.ID == uuid.Nil {
			command.Command = cmd.Command
			command.Group = cmd.Group
		}

		if cmd.Response != nil {
			command.Responses = append(command.Responses, *cmd.Response)
		}
	}

	return command, nil
}
