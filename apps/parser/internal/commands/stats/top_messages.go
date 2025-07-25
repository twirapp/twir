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

var TopMessages = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "top messages",
		Description: null.StringFrom(*top.Messages.Description),
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
					top.Messages.Name,
				),
			},
		}

		return result, nil
	},
}
