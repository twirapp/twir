package stream

import (
	"context"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/apps/parser/pkg/helpers"
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
			result.Result = i18n.GetCtx(ctx, locales.Translations.Variables.Stream.Errors.Offline)
			return &result, nil
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
