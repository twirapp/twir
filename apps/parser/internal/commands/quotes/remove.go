package quotes

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	quotesrepository "github.com/twirapp/twir/libs/repositories/quotes"
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

		channelUUID, err := uuid.Parse(parseCtx.Channel.DBChannelID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotRemove),
				Err:     err,
			}
		}

		err = parseCtx.Services.QuotesRepo.DeleteByChannelIDAndNumber(ctx, channelUUID, number)
		if errors.Is(err, quotesrepository.ErrQuoteNotFound) {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
			return result, nil
		}
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotRemove),
				Err:     err,
			}
		}

		_ = parseCtx.Services.QuotesCacher.Invalidate(ctx, parseCtx.Channel.DBChannelID)

		result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Remove.Removed)}
		return result, nil
	},
}
