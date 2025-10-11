package channel_title

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/samber/lo"
)

const (
	titleHistoryLimitArgName = "limit"
)

var History = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "title history",
		Description: null.StringFrom("Print history of titles."),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MODERATION",
		IsReply:     true,
		Visible:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name:     titleHistoryLimitArgName,
			Min:      lo.ToPtr(1),
			Max:      lo.ToPtr(20),
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		limit := 5
		limitArg := parseCtx.ArgsParser.Get(titleHistoryLimitArgName)
		if limitArg != nil {
			limit = limitArg.Int()
		}

		var histories []*model.ChannelInfoHistory
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Raw(
				`SELECT * FROM (
				SELECT DISTINCT ON (title) * FROM "channels_info_history"
				                             WHERE "channelId" = ?
				                             ORDER BY "title", "createdAt"
				                             DESC
				) subquery ORDER BY "createdAt" DESC LIMIT ?`, parseCtx.Channel.ID, limit,
			).
			Find(&histories).
			Error

		if err != nil {
			return result, fmt.Errorf(i18n.GetCtx(
				ctx,
				locales.Translations.Commands.Channel.Errors.HistoryTitleMessage.
					SetVars(locales.KeysCommandsChannelErrorsHistoryTitleMessageVars{Reason: err.Error()}),
			))
		}

		titles := lo.Map(
			histories, func(item *model.ChannelInfoHistory, _ int) string {
				return item.Title
			},
		)

		result.Result = append(result.Result, strings.Join(titles, " ║ "))
		return result, nil
	},
}
