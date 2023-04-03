package stats

import (
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	user_follow "github.com/satont/tsuwari/apps/parser/internal/variables/user/follow"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var UserFollowage = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "followage",
		Description: null.StringFrom(*user_follow.FollowageVariable.Description),
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
					"user.followage",
				),
			},
		}

		return result
	},
}
