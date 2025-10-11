package stream

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var Viewers = &types.Variable{
	Name:                "stream.viewers",
	Description:         lo.ToPtr("Stream viewers"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		if parseCtx.ChannelStream != nil {
			result.Result = strconv.Itoa(parseCtx.ChannelStream.ViewerCount)
		} else {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Offline)
		}

		return &result, nil
	},
}
