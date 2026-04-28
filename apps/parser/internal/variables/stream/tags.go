package stream

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/shared"
	"github.com/twirapp/twir/apps/parser/locales"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/i18n"
)

var Tags = &types.Variable{
	Name:                "stream.tags",
	Description:         lo.ToPtr("Stream tags"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		if parseCtx.ChannelStream != nil && len(parseCtx.ChannelStream.Tags) > 0 {
			return &types.VariableHandlerResult{Result: strings.Join(parseCtx.ChannelStream.Tags, ", ")}, nil
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

				if channelInfo == nil || len(channelInfo.Stream.CustomTags) == 0 {
					return &types.VariableHandlerResult{Result: i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)}, nil
				}

				return &types.VariableHandlerResult{Result: strings.Join(channelInfo.Stream.CustomTags, ", ")}, nil
			},
		})(ctx, parseCtx, variableData)
	},
}
