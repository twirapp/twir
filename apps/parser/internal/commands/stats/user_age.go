package stats

import (
	"context"
	"fmt"

	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/variables/user"

	"github.com/guregu/null"
	"github.com/lib/pq"

	model "github.com/satont/twir/libs/gomodels"
)

var UserAge = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "age",
		Description: null.StringFrom(*user.Age.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Visible:     true,
		IsReply:     true,
		Aliases:     pq.StringArray{"accountage"},
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
					user.Age.Name,
				),
			},
		}

		return result, nil
	},
}
