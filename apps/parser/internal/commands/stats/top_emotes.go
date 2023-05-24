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

var TopEmotes = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top emotes",
		Description: null.StringFrom(*top.Emotes.Description),
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
					top.Emotes.Name,
				),
			},
		}

		return result
	},
}
