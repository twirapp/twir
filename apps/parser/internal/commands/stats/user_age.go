package stats

import (
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	userage "github.com/satont/tsuwari/apps/parser/internal/variables/user/age"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
)

var UserAge = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "age",
		Description: null.StringFrom(*userage.Variable.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
		Aliases:     pq.StringArray{"accountage"},
	},
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					userage.Variable.Name,
				),
			},
		}

		return result
	},
}
