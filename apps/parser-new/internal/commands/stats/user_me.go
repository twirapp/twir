package stats

import (
	"context"
	"fmt"

	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables/user"

	"github.com/guregu/null"
	"github.com/lib/pq"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var UserMe = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "me",
		Description: null.StringFrom("Prints user statistic."),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Aliases:     pq.StringArray{"stats"},
		Visible:     true,
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s) used emotes · $(%s) watched · $(%s) messages · $(%s) used points · $(%s) songs requested",
					user.Emotes.Name,
					user.Watched.Name,
					user.Messages.Name,
					user.UsedChannelPoints.Name,
					user.SongsRequested.Name,
				),
			},
		}

		return result
	},
}
