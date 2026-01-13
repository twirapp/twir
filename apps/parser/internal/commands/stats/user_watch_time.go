package stats

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/user"
	"github.com/twirapp/twir/apps/parser/locales"

	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var UserWatchTime = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "watchtime",
		Description: null.StringFrom("Prints user watch time."),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Aliases:     pq.StringArray{"watch"},
		Visible:     true,
		Enabled:     false,
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		watching := fmt.Sprintf(
			"$(%s)",
			user.Watched.Name,
		)

		var resultMessage string
		if len(parseCtx.Mentions) > 0 {
			resultMessage = i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Stats.Info.WatchingStreamMentioned.
					SetVars(locales.KeysCommandsStatsInfoWatchingStreamMentionedVars{
						UserName:     parseCtx.Mentions[0].UserName,
						UserWatching: watching,
					}),
			)
		} else {
			resultMessage = i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Stats.Info.WatchingStream.
					SetVars(locales.KeysCommandsStatsInfoWatchingStreamVars{UserWatching: watching}),
			)
		}

		result := &types.CommandsHandlerResult{
			Result: []string{resultMessage},
		}

		return result, nil
	},
}
