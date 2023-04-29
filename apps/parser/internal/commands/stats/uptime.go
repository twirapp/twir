package stats

import (
	"fmt"
	streamuptime "github.com/satont/tsuwari/apps/parser/internal/variables/stream/uptime"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
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
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
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
