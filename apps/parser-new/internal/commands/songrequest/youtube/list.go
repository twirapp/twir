package sr_youtube

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"

	"github.com/samber/lo"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var SrListCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "sr list",
		Description: null.StringFrom("List of requested songs"),
		Visible:     true,
		Module:      "SONGS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}

		url := fmt.Sprintf(
			"%s://%s/p/%s/song-requests",
			lo.If(parseCtx.Services.Config.AppEnv == "development", "http").Else("https"),
			parseCtx.Services.Config.HostName,
			parseCtx.Channel.Name,
		)

		result.Result = append(result.Result, url)
		return result
	},
}
