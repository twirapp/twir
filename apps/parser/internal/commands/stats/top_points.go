package stats

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/internal/variables/top"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var TopPoints = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top points",
		Description: null.StringFrom(*top.ChannelPoints.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					top.ChannelPoints.Name,
				),
			},
		}

		return result
	},
}
