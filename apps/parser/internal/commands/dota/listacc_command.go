package dota

import (
	model "github.com/satont/tsuwari/libs/gomodels"
	"strings"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
)

var ListAccCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "dota listacc",
		Description: lo.ToPtr("List of added dota accounts"),
		RolesNames:  []model.ChannelRoleEnum{model.ChannelRoleTypeBroadcaster},
		Visible:     false,
		Module:      lo.ToPtr("DOTA"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		accounts := GetAccountsByChannelId(ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			result.Result = append(result.Result, NO_ACCOUNTS)
			return result
		}

		result.Result = append(result.Result, strings.Join(*accounts, ", "))
		return result
	},
}
