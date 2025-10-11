package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"
)

var Title = &types.Variable{
	Name:                "stream.title",
	Description:         lo.ToPtr("Stream title"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		if parseCtx.ChannelStream != nil {
			result.Result = parseCtx.ChannelStream.Title
		} else {
			channelInfo := parseCtx.Cacher.GetTwitchChannel(ctx)
			if channelInfo != nil {
				result.Result = channelInfo.Title
			} else {
				result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)
			}
		}

		return &result, nil
	},
}
