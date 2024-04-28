package channel_game

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"

	"github.com/samber/lo"
)

const (
	categoryHistoryLimitArgName = "limit"
)

var History = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "game history",
		Description: null.StringFrom("Print history of games."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name:     categoryHistoryLimitArgName,
			Min:      lo.ToPtr(1),
			Max:      lo.ToPtr(20),
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		limit := 5
		limitArg := parseCtx.ArgsParser.Get(categoryHistoryLimitArgName)
		if limitArg != nil {
			limit = limitArg.Int()
		}

		var histories []*model.ChannelInfoHistory
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Raw(
				`SELECT * FROM (
				SELECT DISTINCT ON (category) * FROM "channels_info_history"
				                             WHERE "channelId" = ?
				                             ORDER BY "category", "createdAt"
				                             DESC
				) subquery ORDER BY "createdAt" DESC LIMIT ?`, parseCtx.Channel.ID, limit,
			).
			Find(&histories).
			Error

		if err != nil {
			result.Result = append(result.Result, "internal error")
			return nil, &types.CommandHandlerError{
				Message: "cannot find used games in database",
				Err:     err,
			}
		}

		categories := lo.Map(
			histories, func(item *model.ChannelInfoHistory, _ int) string {
				return item.Category
			},
		)

		result.Result = append(result.Result, strings.Join(categories, " ║ "))
		return result, nil
	},
}
