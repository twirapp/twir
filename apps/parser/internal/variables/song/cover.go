package song

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var Cover = &types.Variable{
	Name:                "currentsong.imageUrl",
	Description:         lo.ToPtr("Print current song image url."),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		currentSong := parseCtx.Cacher.GetCurrentSong(ctx)
		if currentSong == nil {
			return result, nil
		}

		result.Result = currentSong.Image

		return result, nil
	},
}
