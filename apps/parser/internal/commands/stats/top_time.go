package stats

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/top"

	model "github.com/twirapp/twir/libs/gomodels"
)

var TopTime = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top time",
		Description: null.StringFrom(*top.Watched.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{"top watched"},
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					top.Watched.Name,
				),
			},
		}

		return result, nil
	},
}
