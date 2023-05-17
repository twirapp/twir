package manage

import (
	"context"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"

	"go.uber.org/zap"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var CheckAliasesCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "commands aliases",
		Description: null.StringFrom("Check command aliases"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		if parseCtx.Text == nil {
			result.Result = append(result.Result, "type command name for check aliases.")
			return result
		}

		commandName := strings.ReplaceAll(strings.ToLower(*parseCtx.Text), "!", "")

		cmd := model.ChannelsCommands{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "name" = ?`, parseCtx.Channel.ID, commandName).
			Find(&cmd).
			Error
		if err != nil {
			zap.S().Error(err)
			result.Result = append(result.Result, "internal error")
			return result
		}

		if cmd.ID == "" {
			result.Result = append(result.Result, "command with that name not found.")
			return result
		}

		if len(cmd.Aliases) == 0 {
			result.Result = append(result.Result, "command have no aliases")
			return result
		}

		result.Result = append(result.Result, strings.Join(cmd.Aliases, ", "))
		return result
	},
}
