package quotes

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/jackc/pgx/v5"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	quotesmodel "github.com/twirapp/twir/libs/repositories/quotes/model"
)

const quoteIDArgName = "id"

var Quote = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "quote",
		Description: null.StringFrom("Show a quote"),
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

		channelUUID, err := uuid.Parse(parseCtx.Channel.DBChannelID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotGet),
				Err:     err,
			}
		}

		var quote quotesmodel.Quote

		if parseCtx.ArgsParser.IsExists(quoteIDArgName) {
			number, ok := parseQuoteNumber(parseCtx.ArgsParser.Get(quoteIDArgName).String())
			if !ok {
				result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
				return result, nil
			}

			quote, err = parseCtx.Services.QuotesRepo.GetByChannelIDAndNumber(ctx, channelUUID, number)
			if errors.Is(err, pgx.ErrNoRows) {
				result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.NotFound)}
				return result, nil
			}
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotGet),
					Err:     err,
				}
			}
		} else {
			quote, err = parseCtx.Services.QuotesRepo.GetRandomByChannelID(ctx, channelUUID)
			if errors.Is(err, pgx.ErrNoRows) {
				result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.RandomEmpty)}
				return result, nil
			}
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Quotes.Errors.CannotGet),
					Err:     err,
				}
			}
		}

		result.Result = []string{fmt.Sprintf("#%d %s", quote.Number, quote.Text)}
		return result, nil
	},
}

func parseQuoteNumber(input string) (int, bool) {
	number, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSpace(input), "#"))
	if err != nil || number < 1 {
		return 0, false
	}

	return number, true
}
