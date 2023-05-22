package stats

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/top"

	model "github.com/satont/tsuwari/libs/gomodels"
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
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					top.Watched.Name,
				),
			},
		}

		return result
	},
}
