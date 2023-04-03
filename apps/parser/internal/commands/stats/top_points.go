package stats

import (
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	top_channel_points "github.com/satont/tsuwari/apps/parser/internal/variables/top/channel_points"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var TopPoints = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top points",
		Description: null.StringFrom(*top_channel_points.Variable.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					top_channel_points.Variable.Name,
				),
			},
		}

		return result
	},
}
