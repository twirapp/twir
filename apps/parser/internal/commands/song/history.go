package song

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	currentsong "github.com/satont/twir/apps/parser/internal/variables/song"
	model "github.com/satont/twir/libs/gomodels"
)

var History = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "song history",
		Description: null.StringFrom(*currentsong.History.Description),
		RolesIDS:    pq.StringArray{},
		Module:      "SONGS",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s)",
					currentsong.History.Name,
				),
			},
		}

		return result, nil
	},
}
