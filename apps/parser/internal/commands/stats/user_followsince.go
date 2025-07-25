package stats

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/user"

	model "github.com/twirapp/twir/libs/gomodels"
)

var UserFollowSince = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "followsince",
		Description: null.StringFrom(*user.FollowSince.Description),
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
					user.FollowSince.Name,
				),
			},
		}

		return result, nil
	},
}
