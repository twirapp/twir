package dota2

import (
	"context"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var Wp = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "wp",
		Description: null.StringFrom("Show the streamer's win probability in the active Dota 2 match"),
		Module:      "DOTA",
		Visible:     true,
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if _, err := requireDotaSettings(
			ctx,
			parseCtx,
			func(settings model.ChannelsDotaSettingsCommands) bool { return settings.Wp },
		); err != nil {
			return nil, err
		}

		data, err := getDotaData(ctx, parseCtx)
		if err != nil {
			return nil, err
		}
		if !data.InGame {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.NoActiveMatch),
			}
		}

		probability, available := winProbabilityOutput(data)
		if !available {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.WinProbabilityUnavailable),
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dota.Outputs.WinProbability.SetVars(
						locales.KeysCommandsDotaOutputsWinProbabilityVars{
							Probability: probability,
						},
					),
				),
			},
		}, nil
	},
}
