package stats

import (
	"context"
	"strings"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/user"
	"github.com/twirapp/twir/apps/parser/locales"

	"github.com/guregu/null"
	"github.com/lib/pq"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var UserMe = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "me",
		Description: null.StringFrom("Prints user statistic."),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Aliases:     pq.StringArray{"stats"},
		Visible:     true,
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var slice []string

		slice = append(slice, i18n.GetCtx(
			ctx,
			locales.Translations.Commands.Stats.Info.Watched.
				SetVars(locales.KeysCommandsStatsInfoWatchedVars{UserWatched: user.Watched.Name}),
		))
		slice = append(slice, i18n.GetCtx(
			ctx,
			locales.Translations.Commands.Stats.Info.Messages.
				SetVars(locales.KeysCommandsStatsInfoMessagesVars{UserMessages: user.Messages.Name}),
		))
		slice = append(slice, i18n.GetCtx(
			ctx,
			locales.Translations.Commands.Stats.Info.Emotes.
				SetVars(locales.KeysCommandsStatsInfoEmotesVars{UserEmotes: user.Messages.Name}),
		))
		slice = append(slice, i18n.GetCtx(
			ctx,
			locales.Translations.Commands.Stats.Info.Points.
				SetVars(locales.KeysCommandsStatsInfoPointsVars{UserPoints: user.UsedChannelPoints.Name}),
		))
		slice = append(slice, i18n.GetCtx(
			ctx,
			locales.Translations.Commands.Stats.Info.Songs.
				SetVars(locales.KeysCommandsStatsInfoSongsVars{UserSongs: user.SongsRequested.Name}),
		))

		result := &types.CommandsHandlerResult{
			Result: []string{strings.Join(slice, " · ")},
		}

		return result, nil
	},
}
