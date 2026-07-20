package dota2

import (
	"context"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var Lg = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:    "lg",
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
			func(settings model.ChannelsDotaSettingsCommands) bool { return settings.Lg },
		); err != nil {
			return nil, err
		}

		data, err := getDotaData(ctx, parseCtx)
		if err != nil {
			return nil, err
		}

		lastGame, ok := formatLastGame(data.LastGame)
		if !ok {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.NoLastGame),
			}
		}

		gameResult := i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Outputs.Loss)
		if lastGame.Won {
			gameResult = i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Outputs.Win)
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dota.Outputs.LastGame.SetVars(
						locales.KeysCommandsDotaOutputsLastGameVars{
							Hero:     lastGame.HeroName,
							Kda:      lastGame.KDA,
							Result:   gameResult,
							Duration: lastGame.Duration,
						},
					),
				),
			},
		}, nil
	},
}
