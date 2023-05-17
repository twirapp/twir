package stats

import (
	"context"
	"fmt"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"

	"github.com/guregu/null"
	"github.com/lib/pq"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var Uptime = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "uptime",
		Description: null.StringFrom(*streamuptime.Variable.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
		Aliases:     pq.StringArray{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					streamuptime.Variable.Name,
				),
			},
		}

		return result
	},
}
