package quotes

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	quotesrepository "github.com/twirapp/twir/libs/repositories/quotes"
)

const quoteTextArgName = "text"

var AddQuote = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "quote add",
		Aliases:     []string{"quote +"},
		Description: null.StringFrom("Add a quote"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "MANAGE",
		IsReply:     true,
	},
	Args: []command_arguments.Arg{
		command_arguments.VariadicString{
			Name:     quoteTextArgName,
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (*types.CommandsHandlerResult, error) {
		result := &types.CommandsHandlerResult{}
		if !parseCtx.ArgsParser.IsExists(quoteTextArgName) {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.EmptyText)}
			return result, nil
		}

		text := strings.TrimSpace(parseCtx.ArgsParser.Get(quoteTextArgName).String())
		if text == "" {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.EmptyText)}
			return result, nil
		}

		channelUUID, err := uuid.Parse(parseCtx.Channel.DBChannelID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotAdd),
				Err:     err,
			}
		}

		createInput := quotesrepository.CreateInput{
			ChannelID:   channelUUID,
			Text:        text,
			CreatorID:   lo.ToPtr(parseCtx.Sender.ID),
			CreatorName: lo.ToPtr(parseCtx.Sender.DisplayName),
		}
		if parseCtx.ChannelStream != nil {
			createInput.GameID = lo.ToPtr(parseCtx.ChannelStream.GameId)
			createInput.GameName = lo.ToPtr(parseCtx.ChannelStream.GameName)
		}

		quote, err := parseCtx.Services.QuotesRepo.Create(ctx, createInput)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotAdd),
				Err:     err,
			}
		}

		_ = parseCtx.Services.QuotesCacher.Invalidate(ctx, parseCtx.Channel.DBChannelID)

		result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Add.Added.SetVars(
			locales.KeysCommandsQuotesAddAddedVars{Number: quote.Number},
		))}
		return result, nil
	},
}
