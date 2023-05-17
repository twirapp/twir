package stats

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var TopMessages = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top messages",
		Description: null.StringFrom(*messages.Variable.Description),
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
					messages.Variable.Name,
				),
			},
		}

		return result
	},
}
