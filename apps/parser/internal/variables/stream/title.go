package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/shared"
	"github.com/twirapp/twir/apps/parser/locales"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/i18n"
)

var Title = &types.Variable{
	Name:                "stream.title",
	Description:         lo.ToPtr("Stream title"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		if parseCtx.ChannelStream != nil {
			return &types.VariableHandlerResult{Result: parseCtx.ChannelStream.Title}, nil
		}

		return shared.HandlerByPlatform(map[platformentity.Platform]types.VariableHandler{
			shared.PlatformTwitch: func(
				ctx context.Context,
				parseCtx *types.VariableParseContext,
				variableData *types.VariableData,
			) (*types.VariableHandlerResult, error) {
				result := types.VariableHandlerResult{}

				channelInfo := parseCtx.Cacher.GetTwitchChannel(ctx)
				if channelInfo != nil {
					result.Result = channelInfo.Title
				} else {
					result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)
				}

				return &result, nil
			},
			shared.PlatformKick: func(
				ctx context.Context,
				parseCtx *types.VariableParseContext,
				variableData *types.VariableData,
			) (*types.VariableHandlerResult, error) {
				result := types.VariableHandlerResult{}

				channelInfo, err := shared.GetKickChannel(ctx, parseCtx)
				if err != nil {
					parseCtx.Services.Logger.Sugar().Error(err)
					result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)
					return &result, nil
				}

				if channelInfo != nil {
					result.Result = channelInfo.StreamTitle
				} else {
					result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)
				}

				return &result, nil
			},
		})(ctx, parseCtx, variableData)
	},
}
