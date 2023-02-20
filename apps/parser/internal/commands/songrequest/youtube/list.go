package sr_youtube

import (
	"fmt"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	config "github.com/satont/tsuwari/libs/config"
)

var SrListCommand = types.DefaultCommand{
	Command: types.Command{
		Name:               "sr list",
		Description:        lo.ToPtr("List of requested songs"),
		Permission:         "VIEWER",
		Visible:            false,
		Module:             lo.ToPtr("SONGREQUEST"),
		IsReply:            true,
		KeepResponsesOrder: lo.ToPtr(false),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		cfg := do.MustInvoke[config.Config](di.Provider)

		result := &types.CommandsHandlerResult{}

		url := fmt.Sprintf(
			"%s://%s/p/%s/song-requests",
			lo.If(cfg.AppEnv == "development", "http").Else("https"),
			cfg.HostName,
			ctx.ChannelName,
		)

		result.Result = append(result.Result, url)
		return result
	},
}
