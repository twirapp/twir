package dota2

import (
	"context"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var Gm = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:    "gm",
		Module:  "DOTA",
		IsReply: true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if _, err := requireDotaSettings(
			ctx,
			parseCtx,
			func(settings model.ChannelsDotaSettingsCommands) bool { return settings.Gm },
		); err != nil {
			return nil, err
		}

		data, err := getDotaData(ctx, parseCtx)
		if err != nil {
			return nil, err
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dota.Outputs.Gm.SetVars(
						locales.KeysCommandsDotaOutputsGmVars{
							Medal: medalName(ctx, medalForMMR(data.Mmr)),
						},
					),
				),
			},
		}, nil
	},
}
