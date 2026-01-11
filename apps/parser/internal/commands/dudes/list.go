package dudes

import (
	"context"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/types/types/overlays"
)

var List = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "dudes sprite list",
		Description: null.StringFrom("List of available dudes sprites"),
		RolesIDS:    pq.StringArray{},
		Module:      "DUDES",
		Visible:     true,
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(_ context.Context, _ *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var availableSpritesStr strings.Builder
		for i, v := range overlays.AllDudesSpriteEnumValues {
			availableSpritesStr.WriteString(v.String())
			if i < len(overlays.AllDudesSpriteEnumValues)-1 {
				availableSpritesStr.WriteString(", ")
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{availableSpritesStr.String()},
		}, nil
	},
}
