package stream

import (
	"context"
	"strconv"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/shared"
	"github.com/twirapp/twir/apps/parser/locales"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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
			return &result, nil
		}

		if parseCtx.Platform != shared.PlatformKick {
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Offline)
			return &result, nil
		}

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

				if channelInfo == nil || !channelInfo.Stream.IsLive {
					return &types.VariableHandlerResult{Result: i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Offline)}, nil
				}

				return &types.VariableHandlerResult{Result: strconv.Itoa(channelInfo.Stream.ViewerCount)}, nil
			},
		})(ctx, parseCtx, variableData)
	},
}
