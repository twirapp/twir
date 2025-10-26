package stats

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/user"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/i18n"

	model "github.com/twirapp/twir/libs/gomodels"
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
		vars := []string{
			fmt.Sprintf(
				"(%s) %s",
				user.Watched.Name,
				i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Me.Watched),
			),
			fmt.Sprintf(
				"(%s) %s",
				user.Messages.Name,
				i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Me.Messages),
			),
			fmt.Sprintf(
				"(%s) %s",
				user.Emotes.Name,
				i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Me.Emotes),
			),
			fmt.Sprintf(
				"(%s) %s",
				user.UsedChannelPoints.Name,
				i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Me.Points),
			),
			fmt.Sprintf(
				"$(%s) %s",
				user.SongsRequested.Name,
				i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Me.Songs),
			),
		}

		result := &types.CommandsHandlerResult{
			Result: []string{strings.Join(vars, " · ")},
		}

		return result, nil
	},
}
