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
		Visible:     true,
		Module:      lo.ToPtr("DOTA"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			return []string{NO_ACCOUNTS}
		}

		return []string{strings.Join(*accounts, ", ")}
	},
}
