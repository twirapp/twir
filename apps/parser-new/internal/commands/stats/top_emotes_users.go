package stats

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var TopEmotesUsers = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top emotes users",
		Description: null.StringFrom(*emotes.UsersVariable.Description),
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
					emotes.UsersVariable.Name,
				),
			},
		}

		return result
	},
}
