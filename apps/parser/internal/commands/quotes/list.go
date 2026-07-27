package quotes

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/parser/internal/types"

	model "github.com/twirapp/twir/libs/gomodels"
)

var ListQuotes = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "quotes list",
		Description: null.StringFrom("Link to the public quotes page"),
		Module:      "MANAGE",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}

		link := fmt.Sprintf(
			"%s/p/%s/quotes",
			parseCtx.Services.Config.SiteBaseUrl,
			parseCtx.Channel.Name,
		)

		result.Result = append(result.Result, link)
		return result, nil
	},
}
