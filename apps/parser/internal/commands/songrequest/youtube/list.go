package sr_youtube

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/satont/twir/apps/parser/internal/types"

	model "github.com/satont/twir/libs/gomodels"
)

var SrListCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "sr list",
		Description: null.StringFrom("List of requested songs"),
		Visible:     true,
		Module:      "SONGS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		link := fmt.Sprintf(
			"%s/p/%s/songs-requests",
			parseCtx.Services.Config.SiteBaseUrl,
			parseCtx.Channel.Name,
		)

		result.Result = append(result.Result, link)
		return result, nil
	},
}
