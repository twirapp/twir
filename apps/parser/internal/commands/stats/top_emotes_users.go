package stats

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/variables/top"

	model "github.com/satont/twir/libs/gomodels"
)

var TopEmotesUsers = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top emotes users",
		Description: null.StringFrom(*top.EmotesUsers.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
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
					top.EmotesUsers.Name,
				),
			},
		}

		return result, nil
	},
}
