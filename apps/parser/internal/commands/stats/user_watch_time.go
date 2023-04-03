package stats

import (
	"fmt"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	userwatched "github.com/satont/tsuwari/apps/parser/internal/variables/user/watched"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
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
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{fmt.Sprintf("You watching stream for $(%s)", userwatched.Variable.Name)},
		}

		return result
	},
}
