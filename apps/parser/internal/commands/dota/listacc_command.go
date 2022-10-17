package dota

import (
	"strings"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var ListAccCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "dota listacc",
		Description: lo.ToPtr("List of added dota accounts"),
		Permission:  "BROADCASTER",
		Visible:     false,
		Module:      lo.ToPtr("DOTA"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			result.Result = append(result.Result, NO_ACCOUNTS)
			return result
		}

		result.Result = append(result.Result, strings.Join(*accounts, ", "))
		return result
	},
}
