package manage

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"
)

var CheckAliasesCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases",
		Description: null.StringFrom("Check command aliases"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if parseCtx.Text == nil {
			result.Result = append(result.Result, "you should specify aliase name")
			return result, nil
		}

		commandName := strings.ReplaceAll(strings.ToLower(*parseCtx.Text), "!", "")

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "name" = ?`, parseCtx.Channel.ID, commandName).
			Find(&cmd).
			Error
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot get command",
				Err:     err,
			}
		}

		if cmd.ID == "" {
			result.Result = append(result.Result, "command with that name not found.")
			return result, nil
		}

		if len(cmd.Aliases) == 0 {
			result.Result = append(result.Result, "command have no aliases")
			return result, nil
		}

		result.Result = append(result.Result, strings.Join(cmd.Aliases, ", "))
		return result, nil
	},
}
