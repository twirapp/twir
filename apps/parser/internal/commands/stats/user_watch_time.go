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

var UserWatchTime = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "watchtime",
		Description: null.StringFrom("Prints user watch time."),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Aliases:     pq.StringArray{"watch"},
		Visible:     true,
		Enabled:     false,
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: []string{fmt.Sprintf("You watching stream for $(%s)", user.Watched.Name)},
		}

		return result, nil
	},
}
