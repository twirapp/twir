package stream

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/shared"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/i18n"
)

var Uptime = &types.Variable{
	Name:                "stream.uptime",
	Description:         lo.ToPtr("Prints uptime of stream"),
	CanBeUsedInRegistry: true,
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := types.VariableHandlerResult{}

		if parseCtx.ChannelStream == nil {
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

					if channelInfo == nil || !channelInfo.Stream.IsLive || channelInfo.Stream.StartTime == "" {
						return &types.VariableHandlerResult{Result: i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Offline)}, nil
					}

					startedAt, err := time.Parse(time.RFC3339, channelInfo.Stream.StartTime)
					if err != nil {
						parseCtx.Services.Logger.Sugar().Error(err)
						return &types.VariableHandlerResult{Result: i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Error)}, nil
					}

					return &types.VariableHandlerResult{Result: helpers.Duration(startedAt, &helpers.DurationOpts{UseUtc: true, Hide: helpers.DurationOptsHide{}})}, nil
				},
			})(ctx, parseCtx, variableData)
		}

		result.Result = helpers.Duration(
			parseCtx.ChannelStream.StartedAt, &helpers.DurationOpts{
				UseUtc: true,
				Hide:   helpers.DurationOptsHide{},
			},
		)

		return &result, nil
	},
}
