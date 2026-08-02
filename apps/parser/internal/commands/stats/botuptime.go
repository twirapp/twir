package stats

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	"github.com/twirapp/twir/libs/uptime"
)

var BotUptime = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:               "botuptime",
		Module:             "STATS",
		Visible:            true,
		IsReply:            true,
		RolesIDS:           pq.StringArray{},
		Aliases:            pq.StringArray{},
		KeepResponsesOrder: true,
		Description:        null.StringFrom("Shows bot services uptime, status and external platforms ping."),
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (*types.CommandsHandlerResult, error) {
		statuses, err := uptime.ReadAll(ctx, parseCtx.Services.Redis)
		if err != nil {
			return nil, fmt.Errorf("read uptime statuses: %w", err)
		}

		if len(statuses) == 0 {
			result := []string{trimChatMessage(i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Stats.Botuptime.NoData.SetVars(
					locales.KeysCommandsStatsBotuptimeNoDataVars{
						Services:    i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Botuptime.Services),
						Unavailable: i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Botuptime.Unavailable),
					},
				),
			))}
			result = append(result, trimChatMessage(
				"📡 "+i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Botuptime.Ping)+": "+formatPings(ctx),
			))

			return &types.CommandsHandlerResult{
				Result: result,
			}, nil
		}

		now := time.Now().UTC()
		services, down, restarted := summarizeServices(
			statuses,
			now,
			i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Botuptime.Unavailable),
		)
		result := []string{trimChatMessage("⚙️ " + strings.Join(services, " · "))}

		problems := make([]string, 0, 2)
		if len(down) > 0 {
			problems = append(problems, i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Botuptime.Down)+": "+strings.Join(down, ", "))
		}
		if len(restarted) > 0 {
			problems = append(
				problems,
				i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Botuptime.RecentlyRestarted)+": "+strings.Join(restarted, ", "),
			)
		}
		if len(problems) > 0 {
			result = append(result, trimChatMessage("⚠️ "+strings.Join(problems, " · ")))
		}

		result = append(result, trimChatMessage(
			"📡 "+i18n.GetCtx(ctx, locales.Translations.Commands.Stats.Botuptime.Ping)+": "+formatPings(ctx),
		))

		return &types.CommandsHandlerResult{Result: result}, nil
	},
}
