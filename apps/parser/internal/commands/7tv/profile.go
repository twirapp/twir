package seventv

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	seventvvariables "github.com/satont/twir/apps/parser/internal/variables/7tv"
	model "github.com/satont/twir/libs/gomodels"
)

var Profile = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv profile",
		Description: null.StringFrom("Information about 7tv profile"),
		RolesIDS:    pq.StringArray{},
		Module:      "7tv",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     false,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s) · $(%s) ($(%s) emotes): $(%s)",
					seventvvariables.ProfileLink.Name,
					seventvvariables.EmoteSetName.Name,
					seventvvariables.EmoteSetCount.Name,
					seventvvariables.EmoteSetLink.Name,
				),
			},
		}

		return result, nil
	},
}
