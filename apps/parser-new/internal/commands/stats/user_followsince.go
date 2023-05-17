package stats

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"
)

var UserFollowSince = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "followsince",
		Description: null.StringFrom(*user_follow.FollowsinceVariable.Description),
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
					"user.followsince",
				),
			},
		}

		return result
	},
}
