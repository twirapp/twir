package song

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

var CurrentSong = &types.Variable{
	Name:                "currentsong",
	Description:         lo.ToPtr("Print current played song from Spotify, Last.fm, e.t.c, and also from song requests."),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		currentSong := parseCtx.Cacher.GetCurrentSong(ctx)
		if currentSong == nil {
			return result, nil
		}

		result.Result = currentSong.Name

		return result, nil
	},
}
