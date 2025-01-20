package stats

import (
	"context"
	"fmt"

	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/variables/stream"

	"github.com/guregu/null"
	"github.com/lib/pq"

	model "github.com/satont/twir/libs/gomodels"
)

var Uptime = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "uptime",
		Description: null.StringFrom(*stream.Uptime.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
		Aliases:     pq.StringArray{},
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
					stream.Uptime.Name,
				),
			},
		}

		return result, nil
	},
}
