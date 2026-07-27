package quotes

import (
	"context"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

var RemoveQuote = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "quote remove",
		Aliases:     pq.StringArray{"quote rem", "quote delete", "quote del", "quote -"},
		Description: null.StringFrom("Remove a quote"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name:     quoteIDArgName,
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (*types.CommandsHandlerResult, error) {
		result := &types.CommandsHandlerResult{}
		if !parseCtx.ArgsParser.IsExists(quoteIDArgName) {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
			return result, nil
		}

		number, ok := parseQuoteNumber(parseCtx.ArgsParser.Get(quoteIDArgName).String())
		if !ok {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
			return result, nil
		}

		deleteResult := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND number = ?`, parseCtx.Channel.DBChannelID, number).
			Delete(&model.ChannelsQuotes{})
		if deleteResult.Error != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotRemove),
				Err:     deleteResult.Error,
			}
		}
		if deleteResult.RowsAffected == 0 {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
			return result, nil
		}

		result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Remove.Removed)}
		return result, nil
	},
}
