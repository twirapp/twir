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

var Description = &types.Variable{
	Name:                "stream.description",
	Description:         lo.ToPtr("Stream description"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		return shared.HandlerByPlatform(map[platformentity.Platform]types.VariableHandler{
			shared.PlatformKick: func(
				ctx context.Context,
				parseCtx *types.VariableParseContext,
				variableData *types.VariableData,
			) (*types.VariableHandlerResult, error) {
				channelInfo, err := shared.GetKickChannel(ctx, parseCtx)
				if err != nil {
					parseCtx.Services.Logger.Sugar().Error(err)
					return &types.VariableHandlerResult{Result: i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)}, nil
				}

				if channelInfo == nil || channelInfo.ChannelDescription == "" {
					return &types.VariableHandlerResult{Result: i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)}, nil
				}

				return &types.VariableHandlerResult{Result: channelInfo.ChannelDescription}, nil
			},
		})(ctx, parseCtx, variableData)
	},
}
